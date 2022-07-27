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
	"sync/atomic"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/search"
	"github.com/blugelabs/bluge/search/aggregations"
	"golang.org/x/sync/errgroup"
)

func MultiSearch(ctx context.Context, req bluge.SearchRequest, readers ...*bluge.Reader) (search.DocumentMatchIterator, error) {
	if len(readers) == 0 {
		return nil, nil
	}
	if len(readers) == 1 {
		return readers[0].Search(ctx, req)
	}

	bucketAggs := make(map[string]search.Aggregation)
	bucketAggs["duration"] = aggregations.Duration()

	eg := &errgroup.Group{}
	eg.SetLimit(10)
	docs := make(chan *search.DocumentMatch, len(readers)*2)
	aggs := make(chan *search.Bucket, len(readers))

	docList := &DocumentList{
		bucket: search.NewBucket("", bucketAggs),
	}
	egm := &errgroup.Group{}
	egm.Go(func() error {
		for doc := range docs {
			docList.addDocument(doc)
		}
		return nil
	})
	egm.Go(func() error {
		for agg := range aggs {
			docList.bucket.Merge(agg)
		}
		return nil
	})

	for _, r := range readers {
		r := r
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
	_ = egm.Wait()

	docList.Done()

	return docList, nil
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

func (d *DocumentList) Done() {
	// TODO: sort
	d.bucket.Finish()
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
