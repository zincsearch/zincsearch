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

package v1

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/meta"
)

func TestSearch(t *testing.T) {
	type args struct {
		iQuery *ZincQuery
	}
	tests := []struct {
		name    string
		args    args
		data    []map[string]interface{}
		want    *SearchResponse
		wantErr bool
	}{
		{
			name: "Search Query - Match",
			args: args{
				iQuery: &ZincQuery{
					SearchType: "match",
					Query: QueryParams{
						Term: "Prabhat",
					},
					Source:     false,
					MaxResults: 10,
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - Term",
			args: args{
				iQuery: &ZincQuery{
					SearchType: "term",
					Query: QueryParams{
						Term: "angeles",
					},
					Source:     true,
					MaxResults: 10,
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
				{
					"name": "Leonardo DiCaprio",
					"address": map[string]interface{}{
						"city":  "Los angeles",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - MatchAll",
			args: args{
				iQuery: &ZincQuery{
					SearchType: "matchall",
					Source:     []interface{}{"city"},
					MaxResults: 10,
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - alldocuments",
			args: args{
				iQuery: &ZincQuery{
					SearchType: "alldocuments",
					MaxResults: 10,
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - matchphrase",
			args: args{
				iQuery: &ZincQuery{
					SearchType: "matchphrase",
					Query: QueryParams{
						Term: "San Francisco",
					},
					MaxResults: 10,
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - prefix",
			args: args{
				iQuery: &ZincQuery{
					SearchType: "prefix",
					Query: QueryParams{
						Term: "sa",
					},
					MaxResults: 10,
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco California Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - wildcard",
			args: args{
				iQuery: &ZincQuery{
					SearchType: "wildcard",
					Query: QueryParams{
						Term: "san*",
					},
					MaxResults: 10,
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - fuzzy",
			args: args{
				iQuery: &ZincQuery{
					SearchType: "fuzzy",
					Query: QueryParams{
						Term: "fransisco", // note the wrong spelling
					},
					MaxResults: 10,
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
				{
					"name": "Leonardo DiCaprio",
					"address": map[string]interface{}{
						"city":  "Los angeles",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - querystring",
			args: args{
				iQuery: &ZincQuery{
					SearchType: "querystring",
					Query: QueryParams{
						Term: "angeles",
					},
					MaxResults: 10,
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
				{
					"name": "Leonardo DiCaprio",
					"address": map[string]interface{}{
						"city":  "Los angeles",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - data range",
			args: args{
				iQuery: &ZincQuery{
					SearchType: "daterange",
					Query: QueryParams{
						StartTime: time.Now().UTC().Add(time.Hour * -24),
						EndTime:   time.Now().UTC().Add(time.Hour),
					},
					MaxResults: 10,
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
				{
					"name": "Leonardo DiCaprio",
					"address": map[string]interface{}{
						"city":  "Los angeles",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
		{
			name: "Search Query - aggs",
			args: args{
				iQuery: &ZincQuery{
					SearchType: "matchall",
					MaxResults: 0,
					Aggregations: map[string]AggregationParams{
						"hobby": {
							AggType: "terms",
							Field:   "hobby",
						},
						"time": {
							AggType: "date_range",
							Field:   "@timestamp",
							DateRanges: []AggregationDateRange{{
								From: time.Now().UTC(),
								To:   time.Now().UTC().Add(time.Hour),
							}},
						},
					},
				},
			},
			data: []map[string]interface{}{
				{
					"name": "Prabhat Sharma",
					"address": map[string]interface{}{
						"city":  "San Francisco",
						"state": "California",
					},
					"hobby": "chess",
				},
				{
					"name": "Leonardo DiCaprio",
					"address": map[string]interface{}{
						"city":  "Los angeles",
						"state": "California",
					},
					"hobby": "chess",
				},
			},
		},
	}

	indexName := "Search.index_1"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, err := core.NewIndex(indexName, "disk", nil)
			assert.NoError(t, err)
			assert.NotNil(t, index)
			err = core.StoreIndex(index)
			assert.NoError(t, err)

			if (index.Mappings) == nil {
				index.Mappings = meta.NewMappings()
			}
			index.Mappings.SetProperty("address.city", meta.Property{
				Type:          "text",
				Index:         true,
				Store:         true,
				Highlightable: true,
			})

			for _, d := range tt.data {
				rand.Seed(time.Now().UnixNano())
				docId := rand.Intn(1000)
				err := index.UpdateDocument(strconv.Itoa(docId), d, true)
				assert.NoError(t, err)
			}
			got, err := Search(index, tt.args.iQuery)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, got.Hits.Total.Value, 1)

			err = core.DeleteIndex(indexName)
			assert.NoError(t, err)
		})
	}
}
