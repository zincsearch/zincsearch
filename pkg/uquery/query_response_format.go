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

package uquery

import (
	"github.com/blugelabs/bluge/search"

	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/uquery/aggregation"
)

func FormatResponse(resp *meta.SearchResponse, q *meta.ZincQuery, buckets *search.Bucket) error {
	var err error
	// format aggregations
	if len(q.Aggregations) > 0 {
		resp.Aggregations, err = aggregation.Response(buckets)
		if err != nil {
			return errors.New(errors.ErrorTypeParsingException, err.Error())
		}
		if len(resp.Aggregations) > 0 {
			delete(resp.Aggregations, "count")
			delete(resp.Aggregations, "duration")
			delete(resp.Aggregations, "max_score")
		}
	}

	return nil
}
