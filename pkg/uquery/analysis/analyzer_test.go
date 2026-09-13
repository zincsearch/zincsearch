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

package analysis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vcaesar/riot/analysis"

	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/uquery/analysis/internal/testutil"
)

func TestRequestAnalyzer(t *testing.T) {
	t.Run("nil data", func(t *testing.T) {
		got, err := RequestAnalyzer(nil)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
	t.Run("nil analyzer", func(t *testing.T) {
		got, err := RequestAnalyzer(&meta.IndexAnalysis{})
		assert.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("builtin tokenizer with builtin filters", func(t *testing.T) {
		got, err := RequestAnalyzer(&meta.IndexAnalysis{
			Analyzer: map[string]*meta.Analyzer{
				"a": {
					Tokenizer:   "whitespace",
					CharFilter:  []string{"html_strip"},
					TokenFilter: []string{"lowercase"},
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, []string{"hello", "world"}, testutil.Terms(got["a"].Analyze([]byte("<b>Hello</b> WORLD"))))
	})

	t.Run("custom tokenizer, char_filter and token_filter", func(t *testing.T) {
		got, err := RequestAnalyzer(&meta.IndexAnalysis{
			Analyzer: map[string]*meta.Analyzer{
				"a": {
					Type:        "custom",
					Tokenizer:   "my_tok",
					CharFilter:  []string{"my_char"},
					TokenFilter: []string{"my_filter"},
				},
			},
			Tokenizer: map[string]interface{}{
				"my_tok": map[string]interface{}{"type": "regexp", "pattern": "[a-zA-Z]+"},
			},
			CharFilter: map[string]interface{}{
				"my_char": map[string]interface{}{"type": "mapping", "mappings": []interface{}{"1 => one"}},
			},
			TokenFilter: map[string]interface{}{
				"my_filter": map[string]interface{}{"type": "uppercase"},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"ONE", "FOO"}, testutil.Terms(got["a"].Analyze([]byte("1 foo"))))
	})

	t.Run("filter alias for token_filter", func(t *testing.T) {
		data := &meta.IndexAnalysis{
			Analyzer: map[string]*meta.Analyzer{
				"a": {Tokenizer: "whitespace", Filter: []string{"my_filter"}},
			},
			Filter: map[string]interface{}{
				"my_filter": map[string]interface{}{"type": "uppercase"},
			},
		}
		got, err := RequestAnalyzer(data)
		assert.NoError(t, err)
		assert.Equal(t, []string{"FOO"}, testutil.Terms(got["a"].Analyze([]byte("foo"))))
		// aliases are normalized in place
		assert.Nil(t, data.Filter)
		assert.NotNil(t, data.TokenFilter)
		assert.Nil(t, data.Analyzer["a"].Filter)
		assert.Equal(t, []string{"my_filter"}, data.Analyzer["a"].TokenFilter)
	})

	t.Run("type regexp", func(t *testing.T) {
		got, err := RequestAnalyzer(&meta.IndexAnalysis{
			Analyzer: map[string]*meta.Analyzer{
				"a": {Type: "pattern", Pattern: "[a-z]+", Lowercase: true, Stopwords: []string{"b"}},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "c"}, testutil.Terms(got["a"].Analyze([]byte("a b c"))))
	})
	t.Run("type regexp invalid pattern", func(t *testing.T) {
		got, err := RequestAnalyzer(&meta.IndexAnalysis{
			Analyzer: map[string]*meta.Analyzer{"a": {Type: "regexp", Pattern: "("}},
		})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("type standard with stopwords", func(t *testing.T) {
		got, err := RequestAnalyzer(&meta.IndexAnalysis{
			Analyzer: map[string]*meta.Analyzer{"a": {Type: "standard", Stopwords: []string{"foo"}}},
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"bar"}, testutil.Terms(got["a"].Analyze([]byte("Foo Bar"))))
	})
	t.Run("type stop", func(t *testing.T) {
		got, err := RequestAnalyzer(&meta.IndexAnalysis{
			Analyzer: map[string]*meta.Analyzer{"a": {Type: "STOP", Stopwords: []string{"foo"}}},
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"bar"}, testutil.Terms(got["a"].Analyze([]byte("foo bar"))))
	})
	t.Run("type builtin analyzer with extra filter", func(t *testing.T) {
		got, err := RequestAnalyzer(&meta.IndexAnalysis{
			Analyzer: map[string]*meta.Analyzer{"a": {Type: "whitespace", TokenFilter: []string{"uppercase"}}},
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"FOO", "BAR"}, testutil.Terms(got["a"].Analyze([]byte("foo Bar"))))
	})
	t.Run("type unsupported builtin", func(t *testing.T) {
		got, err := RequestAnalyzer(&meta.IndexAnalysis{
			Analyzer: map[string]*meta.Analyzer{"a": {Type: "nope"}},
		})
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	errCases := []struct {
		name string
		data *meta.IndexAnalysis
	}{
		{
			name: "missing tokenizer and type",
			data: &meta.IndexAnalysis{Analyzer: map[string]*meta.Analyzer{"a": {}}},
		},
		{
			name: "custom type without tokenizer",
			data: &meta.IndexAnalysis{Analyzer: map[string]*meta.Analyzer{"a": {Type: "custom"}}},
		},
		{
			name: "undefined tokenizer",
			data: &meta.IndexAnalysis{Analyzer: map[string]*meta.Analyzer{"a": {Tokenizer: "nope"}}},
		},
		{
			name: "undefined char_filter",
			data: &meta.IndexAnalysis{Analyzer: map[string]*meta.Analyzer{"a": {Tokenizer: "standard", CharFilter: []string{"nope"}}}},
		},
		{
			name: "undefined token_filter",
			data: &meta.IndexAnalysis{Analyzer: map[string]*meta.Analyzer{"a": {Tokenizer: "standard", TokenFilter: []string{"nope"}}}},
		},
		{
			name: "invalid char_filter definition",
			data: &meta.IndexAnalysis{
				Analyzer:   map[string]*meta.Analyzer{"a": {Tokenizer: "standard"}},
				CharFilter: map[string]interface{}{"x": map[string]interface{}{"type": "nope"}},
			},
		},
		{
			name: "invalid token_filter definition",
			data: &meta.IndexAnalysis{
				Analyzer:    map[string]*meta.Analyzer{"a": {Tokenizer: "standard"}},
				TokenFilter: map[string]interface{}{"x": map[string]interface{}{"type": "nope"}},
			},
		},
		{
			name: "invalid tokenizer definition",
			data: &meta.IndexAnalysis{
				Analyzer:  map[string]*meta.Analyzer{"a": {Tokenizer: "standard"}},
				Tokenizer: map[string]interface{}{"x": map[string]interface{}{"type": "nope"}},
			},
		},
	}
	for _, tt := range errCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RequestAnalyzer(tt.data)
			assert.Error(t, err)
			assert.Nil(t, got)
		})
	}
}

func TestQueryAnalyzer(t *testing.T) {
	names := []string{
		"standard", "simple", "keyword", "web", "regexp", "pattern", "stop", "whitespace",
		"gse_standard", "gse_search",
		"ar", "arabic", "cjk", "ckb", "sorani", "da", "danish", "de", "german", "en", "english",
		"es", "spanish", "fa", "persian", "fi", "finnish", "fr", "french", "hi", "hindi",
		"hu", "hungarian", "it", "italian", "nl", "dutch", "no", "norwegian", "pt", "portuguese",
		"ro", "romanian", "ru", "russian", "sv", "swedish", "tr", "turkish",
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			got, err := QueryAnalyzer(nil, name)
			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}

	t.Run("empty name falls back to default", func(t *testing.T) {
		custom := &analysis.Analyzer{}
		got, err := QueryAnalyzer(map[string]*analysis.Analyzer{"default": custom}, "")
		assert.NoError(t, err)
		assert.Same(t, custom, got)
	})
	t.Run("empty name without custom default is unknown", func(t *testing.T) {
		got, err := QueryAnalyzer(nil, "")
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("custom overrides builtin", func(t *testing.T) {
		custom := &analysis.Analyzer{}
		got, err := QueryAnalyzer(map[string]*analysis.Analyzer{"standard": custom}, "standard")
		assert.NoError(t, err)
		assert.Same(t, custom, got)
	})
	t.Run("unknown", func(t *testing.T) {
		got, err := QueryAnalyzer(map[string]*analysis.Analyzer{}, "nope")
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestQueryAnalyzerForField(t *testing.T) {
	custom := &analysis.Analyzer{}
	search := &analysis.Analyzer{}
	data := map[string]*analysis.Analyzer{"my_ana": custom, "my_search": search}

	textProp := meta.NewProperty("text")
	textProp.Analyzer = "my_ana"
	textProp.SearchAnalyzer = "my_search"
	mappings := meta.NewMappings()
	mappings.SetProperty("title", textProp)
	mappings.SetProperty("plain", meta.NewProperty("text"))
	mappings.SetProperty("tag", meta.NewProperty("keyword"))

	t.Run("empty field", func(t *testing.T) {
		a, s := QueryAnalyzerForField(data, mappings, "")
		assert.Nil(t, a)
		assert.Nil(t, s)
	})
	t.Run("non text field", func(t *testing.T) {
		a, s := QueryAnalyzerForField(data, mappings, "tag")
		assert.Nil(t, a)
		assert.Nil(t, s)
	})
	t.Run("text field with analyzers", func(t *testing.T) {
		a, s := QueryAnalyzerForField(data, mappings, "title")
		assert.Same(t, custom, a)
		assert.Same(t, search, s)
	})
	t.Run("text field without analyzers has no default", func(t *testing.T) {
		a, s := QueryAnalyzerForField(data, mappings, "plain")
		assert.Nil(t, a)
		assert.Nil(t, s)
	})
	t.Run("unmapped field uses custom default", func(t *testing.T) {
		def := &analysis.Analyzer{}
		a, s := QueryAnalyzerForField(map[string]*analysis.Analyzer{"default": def}, mappings, "other")
		assert.Same(t, def, a)
		assert.Same(t, def, s)
	})
	t.Run("nil mappings", func(t *testing.T) {
		a, s := QueryAnalyzerForField(nil, nil, "title")
		assert.Nil(t, a)
		assert.Nil(t, s)
	})
}
