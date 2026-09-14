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

	"github.com/zincsearch/zincsearch/pkg/uquery/analysis/internal/testutil"
)

func TestRequestTokenFilter(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		got, err := RequestTokenFilter(nil)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
	t.Run("ok", func(t *testing.T) {
		got, err := RequestTokenFilter(map[string]interface{}{
			"my_lower": map[string]interface{}{"type": "lowercase"},
			"my_trunc": map[string]interface{}{"type": "truncate", "length": float64(2)},
		})
		assert.NoError(t, err)
		assert.Len(t, got, 2)
		assert.Equal(t, []string{"ab"}, testutil.Terms(got["my_trunc"].Filter(testutil.Tokens("abcd"))))
	})
	t.Run("missing type", func(t *testing.T) {
		got, err := RequestTokenFilter(map[string]interface{}{"x": map[string]interface{}{}})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("invalid options", func(t *testing.T) {
		got, err := RequestTokenFilter(map[string]interface{}{"x": map[string]interface{}{"type": "keyword"}})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestRequestTokenFilterSlice(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		got, err := RequestTokenFilterSlice(nil)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
	t.Run("string and object", func(t *testing.T) {
		got, err := RequestTokenFilterSlice([]interface{}{
			"lowercase",
			map[string]interface{}{"type": "truncate", "length": float64(2)},
		})
		assert.NoError(t, err)
		assert.Len(t, got, 2)
	})
	t.Run("object missing type", func(t *testing.T) {
		got, err := RequestTokenFilterSlice([]interface{}{map[string]interface{}{}})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("unknown string", func(t *testing.T) {
		got, err := RequestTokenFilterSlice([]interface{}{"nope"})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("unsupported element type", func(t *testing.T) {
		got, err := RequestTokenFilterSlice([]interface{}{1})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestRequestTokenFilterSingle(t *testing.T) {
	// filters that work without options
	names := []string{
		"apostrophe", "camel_case", "camelcase", "edge_ngram", "elision", "length",
		"lower_case", "lowercase", "ngram", "porter", "stemmer", "reverse", "shingle",
		"trim", "stop", "truncate", "unique", "upper_case", "uppercase", "gse_stop",
		"ar_normalization", "arabic_normalization", "ar_stemmer", "arabic_stemmer",
		"cjk_bigram", "cjk_width",
		"ckb_normalization", "sorani_normalization", "ckb_stemmer", "sorani_stemmer",
		"da_stemmer", "danish_stemmer",
		"de_normalization", "german_normalization", "de_stemmer", "german_stemmer",
		"de_light_stemmer", "german_light_stemmer",
		"en_possessive_stemmer", "english_possessive_stemmer", "en_stemmer", "english_stemmer",
		"es_stemmer", "spanish_stemmer", "es_light_stemmer", "spanish_light_stemmer",
		"fa_normalization", "persian_normalization",
		"fi_stemmer", "finnish_stemmer",
		"fr_elision", "french_elision", "fr_stemmer", "french_stemmer",
		"fr_light_stemmer", "french_light_stemmer", "fr_minimal_stemmer", "french_minimal_stemmer",
		"ga_elision", "irish_elision",
		"hi_normalization", "hindi_normalization", "hi_stemmer", "hindi_stemmer",
		"hu_stemmer", "hungarian_stemmer",
		"in_normalization", "indic_normalization",
		"it_elision", "italian_elision", "it_stemmer", "italian_stemmer",
		"it_light_stemmer", "italian_light_stemmer",
		"nl_stemmer", "dutch_stemmer",
		"no_stemmer", "norwegian_stemmer",
		"pt_light_stemmer", "portuguese_stemmer", "portuguese_light_stemmer",
		"ro_stemmer", "romanian_stemmer",
		"ru_stemmer", "russian_stemmer",
		"sv_stemmer", "swedish_stemmer",
		"tr_stemmer", "turkish_stemmer",
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			got, err := RequestTokenFilterSingle(name, nil)
			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}

	withOptions := []struct {
		name    string
		options interface{}
	}{
		{"dict", map[string]interface{}{"words": []interface{}{"a"}}},
		{"keyword", map[string]interface{}{"keywords": []interface{}{"a"}}},
		{"keyword_marker", map[string]interface{}{"keywords": []interface{}{"a"}}},
		{"regexp", map[string]interface{}{"pattern": "a"}},
		{"pattern_replace", map[string]interface{}{"pattern": "a"}},
		{"unicodenorm", map[string]interface{}{"form": "nfc"}},
	}
	for _, tt := range withOptions {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RequestTokenFilterSingle(tt.name, tt.options)
			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}

	t.Run("case insensitive", func(t *testing.T) {
		got, err := RequestTokenFilterSingle("LowerCase", nil)
		assert.NoError(t, err)
		assert.Equal(t, []string{"abc"}, testutil.Terms(got.Filter(testutil.Tokens("ABC"))))
	})
	t.Run("option error propagates", func(t *testing.T) {
		got, err := RequestTokenFilterSingle("unicodenorm", map[string]interface{}{"form": "bad"})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("unknown", func(t *testing.T) {
		got, err := RequestTokenFilterSingle("nope", nil)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
