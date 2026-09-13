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

package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vcaesar/riot"

	"github.com/zincsearch/zincsearch/pkg/meta"
)

func termOn(field, value string) map[string]interface{} {
	return map[string]interface{}{"term": map[string]interface{}{field: value}}
}

func TestBoolQuery(t *testing.T) {
	mappings := testMappings()

	t.Run("all clauses as objects", func(t *testing.T) {
		q, err := BoolQuery(map[string]interface{}{
			"MUST":     termOn("tag", "a"),
			"should":   termOn("tag", "b"),
			"must_not": termOn("tag", "c"),
			"filter":   termOn("tag", "d"),
		}, mappings, nil)
		assert.NoError(t, err)
		bq := q.(*riot.BooleanQuery)
		assert.Len(t, bq.Musts(), 2) // must + filter wrapper
		assert.Len(t, bq.Shoulds(), 1)
		assert.Len(t, bq.MustNots(), 1)
		assert.Equal(t, 0, bq.MinShould())

		var filter *riot.BooleanQuery
		for _, m := range bq.Musts() {
			if f, ok := m.(*riot.BooleanQuery); ok {
				filter = f
			}
		}
		if assert.NotNil(t, filter) {
			assert.Equal(t, 0.0, filter.Boost())
			assert.Len(t, filter.Musts(), 1)
		}
	})

	t.Run("all clauses as arrays", func(t *testing.T) {
		q, err := BoolQuery(map[string]interface{}{
			"must":     []interface{}{termOn("tag", "a"), termOn("tag", "b")},
			"should":   []interface{}{termOn("tag", "c"), termOn("tag", "d"), termOn("tag", "e")},
			"must_not": []interface{}{termOn("tag", "f")},
			"filter":   []interface{}{termOn("tag", "g"), termOn("tag", "h")},
		}, mappings, nil)
		assert.NoError(t, err)
		bq := q.(*riot.BooleanQuery)
		assert.Len(t, bq.Musts(), 3)
		assert.Len(t, bq.Shoulds(), 3)
		assert.Len(t, bq.MustNots(), 1)
	})

	t.Run("minimum_should_match", func(t *testing.T) {
		for _, msm := range []interface{}{2.0, "2", "50%"} {
			q, err := BoolQuery(map[string]interface{}{
				"should":               []interface{}{termOn("tag", "a"), termOn("tag", "b"), termOn("tag", "c"), termOn("tag", "d")},
				"minimum_should_match": msm,
			}, mappings, nil)
			assert.NoError(t, err)
			assert.Equal(t, 2, q.(*riot.BooleanQuery).MinShould())
		}
		_, err := BoolQuery(map[string]interface{}{"should": termOn("tag", "a"), "minimum_should_match": "bad"}, mappings, nil)
		assert.ErrorContains(t, err, "unsupported MinimumShouldMatch")
	})

	t.Run("empty", func(t *testing.T) {
		q, err := BoolQuery(map[string]interface{}{}, mappings, nil)
		assert.NoError(t, err)
		assert.Empty(t, q.(*riot.BooleanQuery).Musts())
	})

	t.Run("unknown field", func(t *testing.T) {
		_, err := BoolQuery(map[string]interface{}{"nope": termOn("tag", "a")}, mappings, nil)
		assert.ErrorContains(t, err, "unknown field")
	})

	for _, clause := range []string{"must", "should", "must_not", "filter"} {
		t.Run(clause+" bad type", func(t *testing.T) {
			_, err := BoolQuery(map[string]interface{}{clause: "x"}, mappings, nil)
			assert.ErrorContains(t, err, "doesn't support values of type")
		})
		t.Run(clause+" nested error object", func(t *testing.T) {
			_, err := BoolQuery(map[string]interface{}{clause: map[string]interface{}{"nope": map[string]interface{}{}}}, mappings, nil)
			assert.ErrorContains(t, err, "["+clause+"] failed to parse field")
		})
		t.Run(clause+" nested error array", func(t *testing.T) {
			_, err := BoolQuery(map[string]interface{}{clause: []interface{}{map[string]interface{}{"nope": map[string]interface{}{}}}}, mappings, nil)
			assert.ErrorContains(t, err, "["+clause+"] failed to parse field")
		})
	}
}

func TestPrefixWildcardRegexp(t *testing.T) {
	type stringQuery func(map[string]interface{}) (riot.Query, error)
	tests := []struct {
		name  string
		fn    stringQuery
		value func(riot.Query) (string, string, float64)
	}{
		{"prefix", PrefixQuery, func(q riot.Query) (string, string, float64) {
			p := q.(*riot.PrefixQuery)
			return p.Field(), p.Prefix(), p.Boost()
		}},
		{"wildcard", WildcardQuery, func(q riot.Query) (string, string, float64) {
			w := q.(*riot.WildcardQuery)
			return w.Field(), w.Wildcard(), w.Boost()
		}},
		{"regexp", RegexpQuery, func(q riot.Query) (string, string, float64) {
			r := q.(*riot.RegexpQuery)
			return r.Field(), r.Regexp(), r.Boost()
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := tt.fn(map[string]interface{}{"tag": "go*"})
			assert.NoError(t, err)
			field, val, boost := tt.value(q)
			assert.Equal(t, "tag", field)
			assert.Equal(t, "go*", val)
			assert.Equal(t, 1.0, boost)

			q, err = tt.fn(map[string]interface{}{"tag": map[string]interface{}{"VALUE": "z.*", "boost": 2.0, "flags": "ALL", "x": 1}})
			assert.NoError(t, err)
			_, val, boost = tt.value(q)
			assert.Equal(t, "z.*", val)
			assert.Equal(t, 2.0, boost)

			_, err = tt.fn(map[string]interface{}{"a": "x", "b": "y"})
			assert.ErrorContains(t, err, "multiple fields")
			_, err = tt.fn(map[string]interface{}{"tag": 1})
			assert.ErrorContains(t, err, "doesn't support values of type")
		})
	}
}

func TestFuzzyQuery(t *testing.T) {
	mappings := testMappings()

	q, err := FuzzyQuery(map[string]interface{}{"tag": "hello"}, mappings, nil)
	assert.NoError(t, err)
	fq := q.(*riot.FuzzyQuery)
	assert.Equal(t, "tag", fq.Field())
	assert.Equal(t, "hello", fq.Term())
	assert.Equal(t, 1.0, fq.Boost())
	assert.Equal(t, 0, fq.Prefix())

	q, err = FuzzyQuery(map[string]interface{}{"title": map[string]interface{}{"VALUE": "hello", "fuzziness": 2.0, "prefix_length": 1.0, "boost": 3.0, "x": 1}}, mappings, nil)
	assert.NoError(t, err)
	fq = q.(*riot.FuzzyQuery)
	assert.Equal(t, 2, fq.Fuzziness())
	assert.Equal(t, 1, fq.Prefix())
	assert.Equal(t, 3.0, fq.Boost())

	// AUTO on a 2-char term resolves to 0, which leaves riot's default fuzziness (1) untouched
	q, err = FuzzyQuery(map[string]interface{}{"tag": map[string]interface{}{"value": "ab", "fuzziness": "AUTO"}}, mappings, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, q.(*riot.FuzzyQuery).Fuzziness())

	_, err = FuzzyQuery(map[string]interface{}{"a": "x", "b": "y"}, mappings, nil)
	assert.ErrorContains(t, err, "multiple fields")
	_, err = FuzzyQuery(map[string]interface{}{"tag": 1}, mappings, nil)
	assert.ErrorContains(t, err, "doesn't support values of type")
}

func TestParseFuzziness(t *testing.T) {
	tests := []struct {
		name      string
		fuzziness interface{}
		query     string
		want      int
	}{
		{"int", 2.0, "x", 2},
		{"string int", "1", "x", 1},
		{"invalid string", "abc", "x", 0},
		{"auto short", "AUTO", "ab", 0},
		{"auto medium", "auto", "abcd", 1},
		{"auto long", "AUTO", "abcdefg", 2},
		{"auto longest token wins", "AUTO", "a abcdefg", 2},
		{"auto custom", "AUTO:4,8", "abcd", 1},
		{"auto custom below low", "AUTO:4,8", "abc", 0},
		{"auto custom above high", "AUTO:4,8", "abcdefgh", 2},
		{"auto custom invalid low", "AUTO:1,8", "abcdefgh", 0},
		{"auto custom low >= high", "AUTO:8,4", "abcdefgh", 0},
		{"auto malformed keeps defaults", "AUTO:4", "abcdefg", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ParseFuzziness(tt.fuzziness, tt.query, nil))
		})
	}
}

func TestNotImplementedQueries(t *testing.T) {
	tests := map[string]func(map[string]interface{}) (riot.Query, error){
		"boosting":         BoostingQuery,
		"combined_fields":  CombinedFieldsQuery,
		"exists":           ExistsQuery,
		"terms_set":        TermsSetQuery,
		"geo_bounding_box": GeoBoundingBoxQuery,
		"geo_distance":     GeoDistanceQuery,
		"geo_polygon":      GeoPolygonQuery,
		"geo_shape":        GeoShapeQuery,
	}
	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			q, err := fn(map[string]interface{}{})
			assert.Nil(t, q)
			assert.ErrorContains(t, err, "["+name+"] query doesn't support")
		})
	}
}

func TestMatchAllNone(t *testing.T) {
	q, err := MatchAllQuery()
	assert.NoError(t, err)
	assert.IsType(t, &riot.MatchAllQuery{}, q)
	q, err = MatchNoneQuery()
	assert.NoError(t, err)
	assert.IsType(t, &riot.MatchNoneQuery{}, q)
}

func TestQuery(t *testing.T) {
	mappings := testMappings()

	t.Run("nil is match_all", func(t *testing.T) {
		q, err := Query(nil, mappings, nil)
		assert.NoError(t, err)
		assert.IsType(t, &riot.MatchAllQuery{}, q)
	})

	t.Run("meta.Query struct", func(t *testing.T) {
		q, err := Query(&meta.Query{Term: map[string]*meta.TermQuery{"tag": {Value: "go"}}}, mappings, nil)
		assert.NoError(t, err)
		assert.Equal(t, "go", q.(*riot.TermQuery).Term())
	})

	t.Run("dispatch", func(t *testing.T) {
		tests := []struct {
			query map[string]interface{}
			want  riot.Query
		}{
			{map[string]interface{}{"bool": map[string]interface{}{}}, &riot.BooleanQuery{}},
			{map[string]interface{}{"MATCH": map[string]interface{}{"name": "x"}}, &riot.MatchQuery{}},
			{map[string]interface{}{"match_bool_prefix": map[string]interface{}{"name": "x"}}, &riot.BooleanQuery{}},
			{map[string]interface{}{"match_phrase": map[string]interface{}{"name": "x"}}, &riot.MatchPhraseQuery{}},
			{map[string]interface{}{"match_phrase_prefix": map[string]interface{}{"name": "x"}}, &riot.BooleanQuery{}},
			{map[string]interface{}{"multi_match": map[string]interface{}{"query": "x"}}, &riot.BooleanQuery{}},
			{map[string]interface{}{"match_all": map[string]interface{}{}}, &riot.MatchAllQuery{}},
			{map[string]interface{}{"match_none": map[string]interface{}{}}, &riot.MatchNoneQuery{}},
			{map[string]interface{}{"ids": map[string]interface{}{"values": []interface{}{"1"}}}, &riot.BooleanQuery{}},
			{map[string]interface{}{"range": map[string]interface{}{"age": map[string]interface{}{"gte": 1.0}}}, &riot.NumericRangeQuery{}},
			{map[string]interface{}{"regexp": map[string]interface{}{"tag": "a.*"}}, &riot.RegexpQuery{}},
			{map[string]interface{}{"prefix": map[string]interface{}{"tag": "a"}}, &riot.PrefixQuery{}},
			{map[string]interface{}{"fuzzy": map[string]interface{}{"tag": "a"}}, &riot.FuzzyQuery{}},
			{map[string]interface{}{"wildcard": map[string]interface{}{"tag": "a*"}}, &riot.WildcardQuery{}},
			{map[string]interface{}{"term": map[string]interface{}{"tag": "a"}}, &riot.TermQuery{}},
			{map[string]interface{}{"terms": map[string]interface{}{"tag": []interface{}{"a"}}}, &riot.BooleanQuery{}},
		}
		for _, tt := range tests {
			for k := range tt.query {
				t.Run(k, func(t *testing.T) {
					q, err := Query(tt.query, mappings, nil)
					assert.NoError(t, err)
					assert.IsType(t, tt.want, q)
				})
			}
		}
		for _, k := range []string{"query_string", "simple_query_string"} {
			t.Run(k, func(t *testing.T) {
				q, err := Query(map[string]interface{}{k: map[string]interface{}{"query": "hello"}}, mappings, nil)
				assert.NoError(t, err)
				assert.NotNil(t, q)
			})
		}
	})

	t.Run("sub-query errors are wrapped", func(t *testing.T) {
		for _, k := range []string{"bool", "boosting", "match", "match_bool_prefix", "match_phrase", "match_phrase_prefix",
			"combined_fields", "query_string", "simple_query_string", "exists", "ids", "range", "regexp",
			"prefix", "fuzzy", "wildcard", "term", "terms", "terms_set", "geo_bounding_box", "geo_distance", "geo_polygon", "geo_shape"} {
			t.Run(k, func(t *testing.T) {
				// two fields (or a bad child) is rejected by every parser that validates input
				_, err := Query(map[string]interface{}{k: map[string]interface{}{"a": 1, "b": 2, "c": 3}}, mappings, nil)
				assert.ErrorContains(t, err, "["+k+"] failed to parse field")
			})
		}
		// multi_match ignores unknown keys, so trigger its operator validation instead
		_, err := Query(map[string]interface{}{"multi_match": map[string]interface{}{"query": "x", "operator": "xor"}}, mappings, nil)
		assert.ErrorContains(t, err, "[multi_match] failed to parse field")
	})

	t.Run("errors", func(t *testing.T) {
		_, err := Query("x", mappings, nil)
		assert.ErrorContains(t, err, "must be a map")
		_, err = Query(map[string]interface{}{"term": "x"}, mappings, nil)
		assert.ErrorContains(t, err, "doesn't support value type")
		_, err = Query(map[string]interface{}{"nope": map[string]interface{}{}}, mappings, nil)
		assert.ErrorContains(t, err, "[nope] query doesn't support")
		_, err = Query(map[string]interface{}{"match_all": map[string]interface{}{}, "match_none": map[string]interface{}{}}, mappings, nil)
		assert.ErrorContains(t, err, "malformed query")
	})
}
