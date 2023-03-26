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
	"container/heap"
	"context"
	"sync/atomic"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/search"
	"github.com/blugelabs/bluge/search/aggregations"
	"golang.org/x/sync/errgroup"

	"github.com/zinclabs/zincsearch/pkg/config"
	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/pkg/uquery"
)

func MultiSearch(
	ctx context.Context,
	query *meta.ZincQuery,
	mappings *meta.Mappings,
	analyzers map[string]*analysis.Analyzer,
	readers ...*bluge.Reader,
) (search.DocumentMatchIterator, error) {
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
	eg.SetLimit(config.Global.Shard.GorutineNum)
	docs := make(chan *search.DocumentMatch, len(readers)*10)
	aggs := make(chan *search.Bucket, len(readers))

	docList := &DocumentList{
		bucket: search.NewBucket("", bucketAggs),
		from:   int64(query.From),
		size:   int64(query.Size),
	}
	heap.Init(docList)
	// handle skip and limit
	maxSize := int64(query.Size)
	query.Size += query.From
	query.From = 0

	egDoc := &errgroup.Group{}
	egDoc.Go(func() error {
		for doc := range docs {
			heap.Push(docList, &Document{doc})
		}
		return nil
	})
	egDoc.Go(func() error {
		for agg := range aggs {
			docList.bucket.Merge(agg)
		}
		return nil
	})

	for _, r := range readers {
		r := r
		req, err := uquery.ParseQueryDSL(query, mappings, analyzers)
		if err != nil {
			return nil, err
		}
		if docList.sort == nil {
			if req, ok := req.(*bluge.TopNSearch); ok {
				docList.sort = req.SortOrder().Copy()
			}
		}
		eg.Go(func() error {
			var n int64
			dmi, err := r.Search(ctx, req)
			if err != nil {
				return err
			}
			next, err := dmi.Next()
			for err == nil && next != nil {
				n++
				docs <- next
				next, err = dmi.Next()
			}
			aggs <- dmi.Aggregations()

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
	close(aggs)
	_ = egDoc.Wait()

	docList.Done()
	docList.bucket.Aggregation("duration").Finish()

	if docList.size > maxSize {
		docList.size = maxSize
	}

	return docList, nil
}

type Document struct {
	doc *search.DocumentMatch
}

type DocumentList struct {
	from   int64
	size   int64
	len    int64
	next   int64
	docs   []*Document
	bucket *search.Bucket
	sort   search.SortOrder
}

func (d *DocumentList) Done() {
	// do skip
	for i := int64(0); i < d.from && i < int64(d.Len()); i++ {
		heap.Pop(d)
	}
	// log size
	d.len = int64(len(d.docs))
}

func (d *DocumentList) Next() (*search.DocumentMatch, error) {
	if d.next >= d.size || d.next >= d.len {
		return nil, nil
	}
	doc := heap.Pop(d)
	d.next++
	return doc.(*Document).doc, nil
}

func (d *DocumentList) Aggregations() *search.Bucket {
	return d.bucket
}

func (d *DocumentList) Push(doc interface{}) {
	d.docs = append(d.docs, doc.(*Document))
}

func (d *DocumentList) Pop() interface{} {
	n := len(d.docs)
	doc := d.docs[n-1]
	d.docs = d.docs[:n-1]
	return doc
}

func (d *DocumentList) Len() int           { return len(d.docs) }
func (d *DocumentList) Less(i, j int) bool { return d.sort.Compare(d.docs[i].doc, d.docs[j].doc) < 0 }
func (d *DocumentList) Swap(i, j int)      { d.docs[i], d.docs[j] = d.docs[j], d.docs[i] }
