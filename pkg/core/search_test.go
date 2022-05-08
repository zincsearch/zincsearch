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

package core

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "github.com/zinclabs/zinc/pkg/meta/v1"
)

func TestIndex_Search(t *testing.T) {

	type args struct {
		iQuery *v1.ZincQuery
	}
	tests := []struct {
		name    string
		args    args
		data    []map[string]interface{}
		want    *v1.SearchResponse
		wantErr bool
	}{
		{
			name: "Search Query - Match",
			args: args{
				iQuery: &v1.ZincQuery{
					SearchType: "match",
					Query: v1.QueryParams{
						Term: "Prabhat",
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
			},
		},
		{
			name: "Search Query - Term",
			args: args{
				iQuery: &v1.ZincQuery{
					SearchType: "term",
					Query: v1.QueryParams{
						Term: "angeles",
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
		{
			name: "Search Query - MatchAll",
			args: args{
				iQuery: &v1.ZincQuery{
					SearchType: "matchall",
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
			name: "Search Query - wildcard",
			args: args{
				iQuery: &v1.ZincQuery{
					SearchType: "wildcard",
					Query: v1.QueryParams{
						Term: "san*",
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
			},
		},
		{
			name: "Search Query - fuzzy",
			args: args{
				iQuery: &v1.ZincQuery{
					SearchType: "fuzzy",
					Query: v1.QueryParams{
						Term: "fransisco", // note the wrong spelling
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
		{
			name: "Search Query - querystring1",
			args: args{
				iQuery: &v1.ZincQuery{
					SearchType: "querystring",
					Query: v1.QueryParams{
						Term: "angeles",
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(time.Now().UnixNano())
			id := rand.Intn(10000)
			indexName := "Search.index_" + strconv.Itoa(id)

			index, _ := NewIndex(indexName, "disk", UseNewIndexMeta, nil)

			for _, d := range tt.data {
				rand.Seed(time.Now().UnixNano())
				docId := rand.Intn(1000)
				index.UpdateDocument(strconv.Itoa(docId), d, true)
			}

			got, err := index.Search(tt.args.iQuery)
			assert.Nil(t, err)
			assert.Equal(t, 1, got.Hits.Total.Value)
		})
	}

	// os.RemoveAll("data")
}
