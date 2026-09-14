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
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/pkg/meta"
	v1 "github.com/zincsearch/zincsearch/pkg/meta/v1"
)

func mustClause(t *testing.T, q *meta.ZincQuery) []interface{} {
	t.Helper()
	boolQuery, ok := q.Query.(map[string]interface{})["bool"].(map[string]interface{})
	if !assert.True(t, ok) {
		t.FailNow()
	}
	must, ok := boolQuery["must"].([]interface{})
	if !assert.True(t, ok) {
		t.FailNow()
	}
	return must
}

func TestParseQueryDSLFromV1(t *testing.T) {
	t.Run("top level fields", func(t *testing.T) {
		h := &meta.Highlight{}
		q, err := ParseQueryDSLFromV1(&v1.ZincQuery{
			From: 5, MaxResults: 20, Explain: true, Highlight: h, Source: []string{"a"},
			SortFields: []string{"-year", "name"},
		})
		assert.NoError(t, err)
		assert.Equal(t, 5, q.From)
		assert.Equal(t, 20, q.Size)
		assert.True(t, q.Explain)
		assert.Same(t, h, q.Highlight)
		assert.Equal(t, []string{"a"}, q.Source)
		assert.Equal(t, []interface{}{"-year", "name"}, q.Sort)
		assert.Nil(t, q.Aggregations)
	})

	t.Run("no sort fields leaves sort nil", func(t *testing.T) {
		q, err := ParseQueryDSLFromV1(&v1.ZincQuery{})
		assert.NoError(t, err)
		assert.Nil(t, q.Sort)
	})

	t.Run("search types", func(t *testing.T) {
		params := v1.QueryParams{Field: "name", Term: "zinc"}
		tests := []struct {
			searchType string
			want       map[string]interface{}
		}{
			{"alldocuments", map[string]interface{}{"match_all": map[string]interface{}{}}},
			{"matchall", map[string]interface{}{"match_all": map[string]interface{}{}}},
			{"", map[string]interface{}{"match_all": map[string]interface{}{}}},
			{"unknown", map[string]interface{}{"match_all": map[string]interface{}{}}},
			{"wildcard", map[string]interface{}{"wildcard": map[string]interface{}{"name": "zinc"}}},
			{"fuzzy", map[string]interface{}{"fuzzy": map[string]interface{}{"name": "zinc"}}},
			{"term", map[string]interface{}{"term": map[string]interface{}{"name": "zinc"}}},
			{"match", map[string]interface{}{"match": map[string]interface{}{"name": "zinc"}}},
			{"matchphrase", map[string]interface{}{"match_phrase": map[string]interface{}{"name": "zinc"}}},
			{"prefix", map[string]interface{}{"prefix": map[string]interface{}{"name": "zinc"}}},
			{"querystring", map[string]interface{}{"query_string": map[string]interface{}{"query": "zinc"}}},
		}
		for _, tt := range tests {
			t.Run(tt.searchType, func(t *testing.T) {
				q, err := ParseQueryDSLFromV1(&v1.ZincQuery{SearchType: tt.searchType, Query: params})
				assert.NoError(t, err)
				must := mustClause(t, q)
				assert.Equal(t, []interface{}{tt.want}, must)
			})
		}
	})

	t.Run("time range added before search type", func(t *testing.T) {
		start := time.Date(2022, 3, 4, 0, 0, 0, 0, time.UTC)
		end := start.Add(24 * time.Hour)
		q, err := ParseQueryDSLFromV1(&v1.ZincQuery{SearchType: "matchall", Query: v1.QueryParams{StartTime: start, EndTime: end}})
		assert.NoError(t, err)
		must := mustClause(t, q)
		assert.Len(t, must, 2)
		assert.Equal(t, map[string]interface{}{"range": map[string]interface{}{"@timestamp": map[string]interface{}{
			"format": "epoch_millis", "gte": start.UnixMilli(), "lt": end.UnixMilli(),
		}}}, must[0])
		assert.Equal(t, map[string]interface{}{"match_all": map[string]interface{}{}}, must[1])
	})

	t.Run("only start time", func(t *testing.T) {
		start := time.Date(2022, 3, 4, 0, 0, 0, 0, time.UTC)
		q, err := ParseQueryDSLFromV1(&v1.ZincQuery{Query: v1.QueryParams{StartTime: start}})
		assert.NoError(t, err)
		rng := mustClause(t, q)[0].(map[string]interface{})["range"].(map[string]interface{})["@timestamp"].(map[string]interface{})
		assert.Equal(t, start.UnixMilli(), rng["gte"])
		_, hasLT := rng["lt"]
		assert.False(t, hasLT)
	})

	t.Run("daterange on custom field", func(t *testing.T) {
		start := time.Date(2022, 3, 4, 0, 0, 0, 0, time.UTC)
		end := start.Add(time.Hour)
		q, err := ParseQueryDSLFromV1(&v1.ZincQuery{SearchType: "daterange", Query: v1.QueryParams{Field: "created", StartTime: start, EndTime: end}})
		assert.NoError(t, err)
		must := mustClause(t, q)
		assert.Len(t, must, 2)
		assert.Equal(t, map[string]interface{}{"range": map[string]interface{}{"created": map[string]interface{}{
			"gte": start.UnixMilli(), "lt": end.UnixMilli(), "format": "epoch_millis",
		}}}, must[1])
	})

	t.Run("daterange on @timestamp is not duplicated", func(t *testing.T) {
		start := time.Date(2022, 3, 4, 0, 0, 0, 0, time.UTC)
		for _, field := range []string{"", "@timestamp"} {
			q, err := ParseQueryDSLFromV1(&v1.ZincQuery{SearchType: "daterange", Query: v1.QueryParams{Field: field, StartTime: start, EndTime: start.Add(time.Hour)}})
			assert.NoError(t, err)
			assert.Len(t, mustClause(t, q), 1)
		}
	})

	t.Run("aggregations", func(t *testing.T) {
		from := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)
		q, err := ParseQueryDSLFromV1(&v1.ZincQuery{Aggregations: map[string]v1.AggregationParams{
			"t":   {AggType: "terms", Field: "tag", Size: 5},
			"t2":  {AggType: "term", Field: "tag"},
			"r":   {AggType: "range", Field: "age", Ranges: []v1.AggregationNumberRange{{From: 1, To: 10}, {From: 10, To: 20}}},
			"dr":  {AggType: "date_range", Field: "ts", DateRanges: []v1.AggregationDateRange{{From: from, To: to}}},
			"max": {AggType: "max", Field: "age"},
			"min": {AggType: "min", Field: "age"},
			"avg": {AggType: "avg", Field: "age"},
			"wa":  {AggType: "weighted_avg", Field: "age", WeightField: "w"},
			"sum": {AggType: "sum", Field: "age"},
			"cnt": {AggType: "count", Field: "age"},
		}})
		assert.NoError(t, err)
		assert.Len(t, q.Aggregations, 10)
		assert.Equal(t, &meta.AggregationsTerms{Field: "tag", Size: 5}, q.Aggregations["t"].Terms)
		assert.Equal(t, &meta.AggregationsTerms{Field: "tag"}, q.Aggregations["t2"].Terms)
		assert.Equal(t, &meta.AggregationRange{Field: "age", Ranges: []meta.Range{{From: 1, To: 10}, {From: 10, To: 20}}}, q.Aggregations["r"].Range)
		assert.Equal(t, &meta.AggregationDateRange{Field: "ts", Ranges: []meta.DateRange{{From: from.Format(time.RFC3339), To: to.Format(time.RFC3339)}}}, q.Aggregations["dr"].DateRange)
		assert.Equal(t, &meta.AggregationMetric{Field: "age"}, q.Aggregations["max"].Max)
		assert.Equal(t, &meta.AggregationMetric{Field: "age"}, q.Aggregations["min"].Min)
		assert.Equal(t, &meta.AggregationMetric{Field: "age"}, q.Aggregations["avg"].Avg)
		assert.Equal(t, &meta.AggregationMetric{Field: "age", WeightField: "w"}, q.Aggregations["wa"].WeightedAvg)
		assert.Equal(t, &meta.AggregationMetric{Field: "age"}, q.Aggregations["sum"].Sum)
		assert.Equal(t, &meta.AggregationMetric{Field: "age"}, q.Aggregations["cnt"].Count)
	})

	t.Run("nested aggregations", func(t *testing.T) {
		q, err := ParseQueryDSLFromV1(&v1.ZincQuery{Aggregations: map[string]v1.AggregationParams{
			"t": {AggType: "terms", Field: "tag", Aggregations: map[string]v1.AggregationParams{
				"m": {AggType: "max", Field: "age"},
			}},
		}})
		assert.NoError(t, err)
		sub := q.Aggregations["t"].Aggregations
		assert.Len(t, sub, 1)
		assert.Equal(t, &meta.AggregationMetric{Field: "age"}, sub["m"].Max)
	})

	t.Run("unsupported aggregation", func(t *testing.T) {
		_, err := ParseQueryDSLFromV1(&v1.ZincQuery{Aggregations: map[string]v1.AggregationParams{"h": {AggType: "histogram"}}})
		assert.ErrorContains(t, err, "aggregation not supported histogram")

		_, err = ParseQueryDSLFromV1(&v1.ZincQuery{Aggregations: map[string]v1.AggregationParams{
			"t": {AggType: "terms", Field: "tag", Aggregations: map[string]v1.AggregationParams{"h": {AggType: "histogram"}}},
		}})
		assert.ErrorContains(t, err, "aggregation not supported histogram")
	})
}
