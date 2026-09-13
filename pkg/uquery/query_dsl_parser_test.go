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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vcaesar/riot"
	"github.com/vcaesar/riot/search"

	"github.com/zincsearch/zincsearch/pkg/config"
	"github.com/zincsearch/zincsearch/pkg/meta"
)

func dslMappings() *meta.Mappings {
	m := meta.NewMappings()
	m.SetProperty("name", meta.NewProperty("text"))
	m.SetProperty("tag", meta.NewProperty("keyword"))
	m.SetProperty("age", meta.NewProperty("numeric"))
	return m
}

func TestParseQueryDSL(t *testing.T) {
	mappings := dslMappings()
	matchAll := map[string]interface{}{"match_all": map[string]interface{}{}}

	t.Run("defaults", func(t *testing.T) {
		q := &meta.ZincQuery{Query: matchAll, Size: 10}
		req, err := ParseQueryDSL(q, mappings, nil)
		assert.NoError(t, err)
		top := req.(*riot.TopNSearch)
		assert.Equal(t, 10, top.Size())
		assert.Equal(t, 0, top.From())
		// riot seeds a default _score sort when no custom sort is given
		assert.Len(t, top.SortOrder(), 1)
		assert.Nil(t, q.Sort)
		assert.Equal(t, &meta.Source{Enable: true}, q.Source)
		assert.Nil(t, q.Fields)
	})

	t.Run("nil query is match_all", func(t *testing.T) {
		_, err := ParseQueryDSL(&meta.ZincQuery{}, mappings, nil)
		assert.NoError(t, err)
	})

	t.Run("size is capped", func(t *testing.T) {
		max := config.Global.MaxResults
		config.Global.MaxResults = 50
		t.Cleanup(func() { config.Global.MaxResults = max })
		q := &meta.ZincQuery{Query: matchAll, Size: 500}
		req, err := ParseQueryDSL(q, mappings, nil)
		assert.NoError(t, err)
		assert.Equal(t, 50, q.Size)
		assert.Equal(t, 50, req.(*riot.TopNSearch).Size())
	})

	t.Run("from, explain, highlight, fields, source, sort", func(t *testing.T) {
		q := &meta.ZincQuery{
			Query:     map[string]interface{}{"term": map[string]interface{}{"tag": "go"}},
			Size:      10,
			From:      20,
			Explain:   true,
			Highlight: &meta.Highlight{Fields: map[string]*meta.Highlight{"name": {}}},
			Fields:    []interface{}{"name", map[string]interface{}{"field": "age"}},
			Source:    []interface{}{"name"},
			Sort:      []interface{}{"-age", map[string]interface{}{"tag": "asc"}},
		}
		req, err := ParseQueryDSL(q, mappings, nil)
		assert.NoError(t, err)
		top := req.(*riot.TopNSearch)
		assert.Equal(t, 20, top.From())
		assert.Equal(t, 3, q.Highlight.Fields["name"].NumberOfFragments)
		assert.Equal(t, []*meta.Field{{Field: "name"}, {Field: "age"}}, q.Fields)
		assert.Equal(t, &meta.Source{Enable: true, Fields: []string{"name"}}, q.Source)
		sort, ok := q.Sort.(search.SortOrder)
		assert.True(t, ok)
		assert.Len(t, sort, 2)
		assert.Equal(t, []string{"age", "tag"}, top.SortOrder().Fields())
	})

	t.Run("non-array fields left untouched", func(t *testing.T) {
		q := &meta.ZincQuery{Query: matchAll, Size: 1, Fields: "name"}
		_, err := ParseQueryDSL(q, mappings, nil)
		assert.NoError(t, err)
		assert.Equal(t, "name", q.Fields)
	})

	t.Run("aggregations", func(t *testing.T) {
		q := &meta.ZincQuery{Query: matchAll, Size: 1, Aggregations: map[string]meta.Aggregations{
			"tags": {Terms: &meta.AggregationsTerms{Field: "tag"}},
			"avg":  {Avg: &meta.AggregationMetric{Field: "age"}},
		}}
		_, err := ParseQueryDSL(q, mappings, nil)
		assert.NoError(t, err)

		q = &meta.ZincQuery{Query: matchAll, Size: 1, Aggregations: map[string]meta.Aggregations{
			"bad": {Terms: &meta.AggregationsTerms{Field: "missing"}},
		}}
		_, err = ParseQueryDSL(q, mappings, nil)
		assert.ErrorContains(t, err, "[terms] aggregation doesn't support")
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name  string
			query *meta.ZincQuery
			want  string
		}{
			{"unknown query", &meta.ZincQuery{Query: map[string]interface{}{"nope": map[string]interface{}{}}}, "[nope] query doesn't support"},
			{"query not a map", &meta.ZincQuery{Query: "x"}, "must be a map"},
			{"bad fields", &meta.ZincQuery{Query: matchAll, Fields: []interface{}{1}}, "[fields]"},
			{"bad source", &meta.ZincQuery{Query: matchAll, Source: 1}, "[_source]"},
			{"bad sort", &meta.ZincQuery{Query: matchAll, Sort: 1}, "[sort]"},
			{"bad sort object", &meta.ZincQuery{Query: matchAll, Sort: []interface{}{map[string]interface{}{"a": "asc", "b": "asc"}}}, "[sort]"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req, err := ParseQueryDSL(tt.query, mappings, nil)
				assert.Nil(t, req)
				assert.ErrorContains(t, err, tt.want)
			})
		}
	})
}
