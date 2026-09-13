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
)

const (
	opOr  = riot.MatchQueryOperator(riot.MatchQueryOperatorOr)
	opAnd = riot.MatchQueryOperator(riot.MatchQueryOperatorAnd)
)

func TestMatchQuery(t *testing.T) {
	mappings := testMappings()

	t.Run("string", func(t *testing.T) {
		q, err := MatchQuery(map[string]interface{}{"name": "hello world"}, mappings, nil)
		assert.NoError(t, err)
		mq := q.(*riot.MatchQuery)
		assert.Equal(t, "name", mq.Field())
		assert.Equal(t, "hello world", mq.Match())
		// text field without a configured analyzer resolves to none
		assert.Nil(t, mq.Analyzer())
		assert.Equal(t, 1.0, mq.Boost())
		assert.Equal(t, opOr, mq.Operator())
	})

	t.Run("field analyzer from mappings", func(t *testing.T) {
		q, err := MatchQuery(map[string]interface{}{"title": "hello"}, mappings, nil)
		assert.NoError(t, err)
		assert.NotNil(t, q.(*riot.MatchQuery).Analyzer())
	})

	t.Run("object with options", func(t *testing.T) {
		q, err := MatchQuery(map[string]interface{}{"name": map[string]interface{}{
			"QUERY": "hello world", "operator": "and", "fuzziness": 1.0, "prefix_length": 2.0, "boost": 2.0, "analyzer": "keyword", "unknown": 1,
		}}, mappings, nil)
		assert.NoError(t, err)
		mq := q.(*riot.MatchQuery)
		assert.Equal(t, opAnd, mq.Operator())
		assert.Equal(t, 1, mq.Fuzziness())
		assert.Equal(t, 2, mq.Prefix())
		assert.Equal(t, 2.0, mq.Boost())
		assert.NotNil(t, mq.Analyzer())
	})

	t.Run("or operator", func(t *testing.T) {
		q, err := MatchQuery(map[string]interface{}{"name": map[string]interface{}{"query": "x", "operator": "OR"}}, mappings, nil)
		assert.NoError(t, err)
		assert.Equal(t, opOr, q.(*riot.MatchQuery).Operator())
	})

	t.Run("auto fuzziness", func(t *testing.T) {
		q, err := MatchQuery(map[string]interface{}{"name": map[string]interface{}{"query": "elasticsearch", "fuzziness": "AUTO"}}, mappings, nil)
		assert.NoError(t, err)
		assert.Equal(t, 2, q.(*riot.MatchQuery).Fuzziness())
	})

	t.Run("keyword field has no analyzer", func(t *testing.T) {
		q, err := MatchQuery(map[string]interface{}{"tag": "x"}, mappings, nil)
		assert.NoError(t, err)
		assert.Nil(t, q.(*riot.MatchQuery).Analyzer())
	})

	t.Run("minimum_should_match builds boolean", func(t *testing.T) {
		q, err := MatchQuery(map[string]interface{}{"name": map[string]interface{}{"query": "a b c", "minimum_should_match": 2.0, "boost": 1.5}}, mappings, nil)
		assert.NoError(t, err)
		bq := q.(*riot.BooleanQuery)
		assert.Len(t, bq.Shoulds(), 3)
		assert.Equal(t, 2, bq.MinShould())
		assert.Equal(t, 1.5, bq.Boost())
		assert.IsType(t, &riot.TermQuery{}, bq.Shoulds()[0])
		assert.Equal(t, 1.5, bq.Shoulds()[0].(*riot.TermQuery).Boost())
	})

	t.Run("minimum_should_match with fuzziness uses fuzzy terms", func(t *testing.T) {
		q, err := MatchQuery(map[string]interface{}{"tag": map[string]interface{}{"query": "hello world", "minimum_should_match": "50%", "fuzziness": 1.0, "prefix_length": 1.0}}, mappings, nil)
		assert.NoError(t, err)
		bq := q.(*riot.BooleanQuery)
		assert.Len(t, bq.Shoulds(), 2)
		assert.Equal(t, 1, bq.MinShould())
		fq := bq.Shoulds()[0].(*riot.FuzzyQuery)
		assert.Equal(t, 1, fq.Fuzziness())
		assert.Equal(t, 1, fq.Prefix())
		assert.Equal(t, "tag", fq.Field())
	})

	t.Run("minimum_should_match ignored with AND", func(t *testing.T) {
		q, err := MatchQuery(map[string]interface{}{"name": map[string]interface{}{"query": "a b", "minimum_should_match": 2.0, "operator": "and"}}, mappings, nil)
		assert.NoError(t, err)
		assert.IsType(t, &riot.MatchQuery{}, q)
	})

	t.Run("minimum_should_match empty query is match_none", func(t *testing.T) {
		q, err := MatchQuery(map[string]interface{}{"name": map[string]interface{}{"query": "", "minimum_should_match": 1.0}}, mappings, nil)
		assert.NoError(t, err)
		assert.IsType(t, &riot.MatchNoneQuery{}, q)
	})

	t.Run("minimum_should_match invalid", func(t *testing.T) {
		_, err := MatchQuery(map[string]interface{}{"name": map[string]interface{}{"query": "a b", "minimum_should_match": "bad"}}, mappings, nil)
		assert.Error(t, err)
	})

	t.Run("errors", func(t *testing.T) {
		_, err := MatchQuery(map[string]interface{}{"a": "x", "b": "y"}, mappings, nil)
		assert.ErrorContains(t, err, "multiple fields")
		_, err = MatchQuery(map[string]interface{}{"name": 1}, mappings, nil)
		assert.ErrorContains(t, err, "doesn't support values of type")
		_, err = MatchQuery(map[string]interface{}{"name": map[string]interface{}{"query": "x", "operator": "xor"}}, mappings, nil)
		assert.ErrorContains(t, err, "unknown operator")
		_, err = MatchQuery(map[string]interface{}{"name": map[string]interface{}{"query": "x", "analyzer": "nope"}}, mappings, nil)
		assert.ErrorContains(t, err, "unknown analyzer")
	})
}

func TestMatchPhraseQuery(t *testing.T) {
	mappings := testMappings()

	q, err := MatchPhraseQuery(map[string]interface{}{"title": "hello world"}, mappings, nil)
	assert.NoError(t, err)
	pq := q.(*riot.MatchPhraseQuery)
	assert.Equal(t, "title", pq.Field())
	assert.Equal(t, "hello world", pq.Phrase())
	assert.NotNil(t, pq.Analyzer())
	assert.Equal(t, 1.0, pq.Boost())

	q, err = MatchPhraseQuery(map[string]interface{}{"tag": map[string]interface{}{"Query": "a b", "analyzer": "simple", "boost": 2.0, "x": 1}}, mappings, nil)
	assert.NoError(t, err)
	pq = q.(*riot.MatchPhraseQuery)
	assert.Equal(t, "a b", pq.Phrase())
	assert.NotNil(t, pq.Analyzer())
	assert.Equal(t, 2.0, pq.Boost())

	q, err = MatchPhraseQuery(map[string]interface{}{"tag": "a"}, mappings, nil)
	assert.NoError(t, err)
	assert.Nil(t, q.(*riot.MatchPhraseQuery).Analyzer())

	_, err = MatchPhraseQuery(map[string]interface{}{"a": "x", "b": "y"}, mappings, nil)
	assert.ErrorContains(t, err, "multiple fields")
	_, err = MatchPhraseQuery(map[string]interface{}{"name": 1}, mappings, nil)
	assert.ErrorContains(t, err, "doesn't support values of type")
	_, err = MatchPhraseQuery(map[string]interface{}{"name": map[string]interface{}{"query": "x", "analyzer": "nope"}}, mappings, nil)
	assert.ErrorContains(t, err, "unknown analyzer")
}

func TestMatchBoolPrefixQuery(t *testing.T) {
	mappings := testMappings()

	q, err := MatchBoolPrefixQuery(map[string]interface{}{"name": "quick brown f"}, mappings, nil)
	assert.NoError(t, err)
	bq := q.(*riot.BooleanQuery)
	assert.Len(t, bq.Shoulds(), 3)
	assert.Equal(t, "quick", bq.Shoulds()[0].(*riot.TermQuery).Term())
	assert.Equal(t, "brown", bq.Shoulds()[1].(*riot.TermQuery).Term())
	assert.Equal(t, "f", bq.Shoulds()[2].(*riot.PrefixQuery).Prefix())
	assert.Equal(t, "name", bq.Shoulds()[2].(*riot.PrefixQuery).Field())
	assert.Equal(t, 1.0, bq.Boost())

	// keyword field falls back to standard analyzer; explicit analyzer + boost honored
	q, err = MatchBoolPrefixQuery(map[string]interface{}{"tag": map[string]interface{}{"query": "Hello", "analyzer": "keyword", "boost": 2.0, "x": 1}}, mappings, nil)
	assert.NoError(t, err)
	bq = q.(*riot.BooleanQuery)
	assert.Len(t, bq.Shoulds(), 1)
	assert.Equal(t, "Hello", bq.Shoulds()[0].(*riot.PrefixQuery).Prefix())
	assert.Equal(t, 2.0, bq.Boost())

	q, err = MatchBoolPrefixQuery(map[string]interface{}{"tag": "a b"}, mappings, nil)
	assert.NoError(t, err)
	assert.Len(t, q.(*riot.BooleanQuery).Shoulds(), 2)

	q, err = MatchBoolPrefixQuery(map[string]interface{}{"name": ""}, mappings, nil)
	assert.NoError(t, err)
	assert.Empty(t, q.(*riot.BooleanQuery).Shoulds())

	_, err = MatchBoolPrefixQuery(map[string]interface{}{"a": "x", "b": "y"}, mappings, nil)
	assert.ErrorContains(t, err, "multiple fields")
	_, err = MatchBoolPrefixQuery(map[string]interface{}{"name": 1}, mappings, nil)
	assert.ErrorContains(t, err, "doesn't support values of type")
	_, err = MatchBoolPrefixQuery(map[string]interface{}{"name": map[string]interface{}{"query": "x", "analyzer": "nope"}}, mappings, nil)
	assert.ErrorContains(t, err, "unknown analyzer")
}

func TestMatchPhrasePrefixQuery(t *testing.T) {
	mappings := testMappings()

	q, err := MatchPhrasePrefixQuery(map[string]interface{}{"name": "quick brown f"}, mappings, nil)
	assert.NoError(t, err)
	bq := q.(*riot.BooleanQuery)
	assert.Len(t, bq.Musts(), 2)
	assert.Equal(t, "f", bq.Musts()[0].(*riot.PrefixQuery).Prefix())
	assert.Equal(t, "quick brown", bq.Musts()[1].(*riot.MatchPhraseQuery).Phrase())
	assert.Equal(t, "name", bq.Musts()[1].(*riot.MatchPhraseQuery).Field())
	assert.Equal(t, 1.0, bq.Boost())

	q, err = MatchPhrasePrefixQuery(map[string]interface{}{"tag": map[string]interface{}{"query": "one", "analyzer": "keyword", "boost": 2.0, "x": 1}}, mappings, nil)
	assert.NoError(t, err)
	bq = q.(*riot.BooleanQuery)
	assert.Len(t, bq.Musts(), 1)
	assert.Equal(t, "one", bq.Musts()[0].(*riot.PrefixQuery).Prefix())
	assert.Equal(t, 2.0, bq.Boost())

	q, err = MatchPhrasePrefixQuery(map[string]interface{}{"tag": ""}, mappings, nil)
	assert.NoError(t, err)
	assert.Empty(t, q.(*riot.BooleanQuery).Musts())

	_, err = MatchPhrasePrefixQuery(map[string]interface{}{"a": "x", "b": "y"}, mappings, nil)
	assert.ErrorContains(t, err, "multiple fields")
	_, err = MatchPhrasePrefixQuery(map[string]interface{}{"name": 1}, mappings, nil)
	assert.ErrorContains(t, err, "doesn't support values of type")
	_, err = MatchPhrasePrefixQuery(map[string]interface{}{"name": map[string]interface{}{"query": "x", "analyzer": "nope"}}, mappings, nil)
	assert.ErrorContains(t, err, "unknown analyzer")
}

func TestMultiMatchQuery(t *testing.T) {
	mappings := testMappings()

	t.Run("basic", func(t *testing.T) {
		q, err := MultiMatchQuery(map[string]interface{}{"QUERY": "hello", "fields": []interface{}{"name", "tag"}, "type": "best_fields", "x": 1}, mappings, nil)
		assert.NoError(t, err)
		bq := q.(*riot.BooleanQuery)
		assert.Len(t, bq.Shoulds(), 2)
		assert.Equal(t, 0, bq.MinShould())
		assert.Equal(t, 1.0, bq.Boost())
		for i, field := range []string{"name", "tag"} {
			mq := bq.Shoulds()[i].(*riot.MatchQuery)
			assert.Equal(t, field, mq.Field())
			assert.Equal(t, "hello", mq.Match())
			assert.Equal(t, opOr, mq.Operator())
		}
	})

	t.Run("options", func(t *testing.T) {
		q, err := MultiMatchQuery(map[string]interface{}{"query": "hello", "fields": []interface{}{"name", "tag", "other"},
			"operator": "AND", "boost": 2.0, "minimum_should_match": 2.0, "analyzer": "keyword"}, mappings, nil)
		assert.NoError(t, err)
		bq := q.(*riot.BooleanQuery)
		assert.Len(t, bq.Shoulds(), 3)
		assert.Equal(t, 2, bq.MinShould())
		assert.Equal(t, 2.0, bq.Boost())
		assert.Equal(t, opAnd, bq.Shoulds()[0].(*riot.MatchQuery).Operator())
		assert.NotNil(t, bq.Shoulds()[0].(*riot.MatchQuery).Analyzer())
	})

	t.Run("or operator and no fields", func(t *testing.T) {
		q, err := MultiMatchQuery(map[string]interface{}{"query": "hello", "operator": "or"}, mappings, nil)
		assert.NoError(t, err)
		assert.Empty(t, q.(*riot.BooleanQuery).Shoulds())
	})

	t.Run("errors", func(t *testing.T) {
		_, err := MultiMatchQuery(map[string]interface{}{"query": "x", "operator": "xor"}, mappings, nil)
		assert.ErrorContains(t, err, "unknown operator")
		_, err = MultiMatchQuery(map[string]interface{}{"query": "x", "fields": []interface{}{"a"}, "minimum_should_match": "bad"}, mappings, nil)
		assert.ErrorContains(t, err, "unsupported MinimumShouldMatch")
	})
}

func TestQueryStringQuery(t *testing.T) {
	mappings := testMappings()

	t.Run("default fields from mappings", func(t *testing.T) {
		q, err := QueryStringQuery(map[string]interface{}{"query": "name:hello AND tag:go"}, mappings, nil)
		assert.NoError(t, err)
		assert.NotNil(t, q)
	})

	t.Run("explicit options", func(t *testing.T) {
		q, err := QueryStringQuery(map[string]interface{}{
			"QUERY": "hello", "fields": []interface{}{"name", "tag", "missing"}, "analyzer": "standard",
			"default_field": "name", "default_operator": "AND", "boost": 2.0, "analyze_wildcard": true,
		}, mappings, nil)
		assert.NoError(t, err)
		assert.NotNil(t, q)
	})

	t.Run("simple_query_string delegates", func(t *testing.T) {
		q, err := SimpleQueryStringQuery(map[string]interface{}{"query": "hello"}, mappings, nil)
		assert.NoError(t, err)
		assert.NotNil(t, q)
	})

	t.Run("unsupported child", func(t *testing.T) {
		_, err := QueryStringQuery(map[string]interface{}{"query": "x", "lenient": true}, mappings, nil)
		assert.ErrorContains(t, err, "unsupported children")
	})
}
