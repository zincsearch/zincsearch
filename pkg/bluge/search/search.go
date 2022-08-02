/* Copyright 2022 Zinc Labs Inc. and Contributors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package search

import (
	"context"
	"github.com/zinclabs/zinc/pkg/config"
	"sort"
	"sync/atomic"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/search"
	"github.com/blugelabs/bluge/search/aggregations"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/uquery"
	"golang.org/x/sync/errgroup"
)

func MultiSearch(ctx context.Context, query *meta.ZincQuery, mappings *meta.Mappings, analyzers map[string]*analysis.Analyzer, readers ...*bluge.Reader) (search.DocumentMatchIterator, error) {
	if len(readers) == 0 {
		return &DocumentList{
			bucket: search.NewBucket("",
				map[string]search.Aggregation{
					"duration": aggregations.Duration(),
				},
			),
		}, nil
	}
	if len(readers) == 1 {
		req, err := uquery.ParseQueryDSL(query, mappings, analyzers)
		if err != nil {
			return nil, err
		}
		return readers[0].Search(ctx, req)
	}

	bucketAggs := make(map[string]search.Aggregation)
	bucketAggs["duration"] = aggregations.Duration()

	eg := &errgroup.Group{}
	eg.SetLimit(config.Global.ReadGorutineNum)
	docs := make(chan *search.DocumentMatch, len(readers)*10)
	aggsChan := make(chan *search.Bucket, len(readers))

	docList := &DocumentList{}
	egm := &errgroup.Group{}
	egm.Go(func() error {
		hitNum := 0
		for doc := range docs {
			hitNum++
			doc.HitNumber = hitNum
			docList.bucket.Consume(doc)
			docList.addDocument(doc)
		}
		docList.bucket.Finish()
		return nil
	})

	var (
		sortOrder    search.SortOrder
		size         int
		skip         int
		reversed     bool
		aggs         search.Aggregations
		neededFields []string
	)

	for _, r := range readers {
		req, err := uquery.ParseQueryDSL(query, mappings, analyzers)
		if err != nil {
			return nil, err
		}
		if sortOrder == nil { // init vars
			aggs = req.Aggregations()
			docList.bucket = search.NewBucket("", aggs)
			sortOrder = req.SortOrder().Copy()
			size, skip, reversed = req.SizeSkipAndReversed()

			neededFields = sortOrder.Fields()
			neededFields = append(neededFields, aggs.Fields()...)
			neededFields = filterRepeatedFields(neededFields)
		}

		r := r
		eg.Go(func() error {
			var n int64
			searcher, err := r.Searcher(req)
			if err != nil {
				return err
			}

			sctx := search.NewSearchContext(size+searcher.DocumentMatchPoolSize(), len(sortOrder))

			next, err := searcher.Next(sctx)
			for err == nil && next != nil {
				n++

				if len(neededFields) > 0 {
					err = next.LoadDocumentValues(sctx, neededFields)
					if err != nil {
						return err
					}
				}

				req.SortOrder().Compute(next)
				docs <- next
				next, err = searcher.Next(sctx)
			}

			if n > atomic.LoadInt64(&docList.size) {
				atomic.StoreInt64(&docList.size, n)
			}

			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	close(docs)
	close(aggsChan)
	_ = egm.Wait()

	err := docList.Done(size, skip, reversed, sortOrder)
	if err != nil {
		return nil, err
	}

	return docList, nil
}

func filterRepeatedFields(s []string) []string {
	if len(s) > 1 {
		filtered := s[:0] // reuse backing array
		store := make(map[string]struct{}, len(s))
		for _, field := range s {
			store[field] = struct{}{}
		}

		for field := range store {
			filtered = append(filtered, field)
		}

		return filtered
	}
	return s
}

type DocumentList struct {
	docs   []*search.DocumentMatch
	bucket *search.Bucket
	size   int64
	next   int64
}

func (d *DocumentList) addDocument(doc *search.DocumentMatch) {
	d.docs = append(d.docs, doc)
}

func (d *DocumentList) Done(size, skip int, reversed bool, sortOrder search.SortOrder) error {
	sort.SliceStable(d.docs, func(i, j int) bool {
		cmp := sortOrder.Compare(d.docs[i], d.docs[j])
		return cmp > 0
	})

	if len(d.docs) > skip {
		d.docs = d.docs[skip:]
	} else {
		d.docs = d.docs[:0]
	}

	// cut down to desired size
	if len(d.docs) > size {
		d.docs = d.docs[:size]
	}

	// complete only selected docs
	for i := range d.docs {
		d.docs[i].Complete(nil)
	}

	if reversed {
		for i, j := 0, len(d.docs)-1; i < j; i, j = i+1, j-1 {
			d.docs[i], d.docs[j] = d.docs[j], d.docs[i]
		}
	}

	return nil
}

func (d *DocumentList) Next() (*search.DocumentMatch, error) {
	if d.next >= d.size || d.next >= int64(len(d.docs)) {
		return nil, nil
	}
	doc := d.docs[d.next]
	d.next++
	return doc, nil
}

func (d *DocumentList) Aggregations() *search.Bucket {
	return d.bucket
}
