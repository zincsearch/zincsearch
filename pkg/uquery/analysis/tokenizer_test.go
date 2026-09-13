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
)

func TestRequestTokenizer(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		got, err := RequestTokenizer(nil)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
	t.Run("ok", func(t *testing.T) {
		got, err := RequestTokenizer(map[string]interface{}{
			"my_ws": map[string]interface{}{"type": "whitespace"},
			"my_re": map[string]interface{}{"type": "regexp", "pattern": "[a-z]+"},
		})
		assert.NoError(t, err)
		assert.Len(t, got, 2)
		assert.Equal(t, []string{"a", "b"}, terms(got["my_re"].Tokenize([]byte("a1b"))))
	})
	t.Run("missing type", func(t *testing.T) {
		got, err := RequestTokenizer(map[string]interface{}{"x": map[string]interface{}{}})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("invalid options", func(t *testing.T) {
		got, err := RequestTokenizer(map[string]interface{}{"x": map[string]interface{}{"type": "char_group"}})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestRequestTokenizerSlice(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		got, err := RequestTokenizerSlice(nil)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
	t.Run("ok", func(t *testing.T) {
		got, err := RequestTokenizerSlice([]interface{}{"whitespace", "standard"})
		assert.NoError(t, err)
		assert.Len(t, got, 2)
	})
	t.Run("not string", func(t *testing.T) {
		got, err := RequestTokenizerSlice([]interface{}{map[string]interface{}{"type": "whitespace"}})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("unknown", func(t *testing.T) {
		got, err := RequestTokenizerSlice([]interface{}{"nope"})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestRequestTokenizerSingle(t *testing.T) {
	names := []string{
		"edge_ngram", "letter", "simple", "lower_case", "lowercase", "ngram", "path_hierarchy",
		"regexp", "pattern", "single", "keyword", "unicode", "standard", "web", "whitespace",
		"gse_standard", "gse_search",
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			got, err := RequestTokenizerSingle(name, nil)
			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}

	withOptions := []struct {
		name    string
		options interface{}
	}{
		{"character", map[string]interface{}{"char": "letter"}},
		{"char_group", map[string]interface{}{"tokenize_on_chars": []interface{}{"-"}}},
		{"exception", map[string]interface{}{"patterns": []interface{}{"a"}}},
	}
	for _, tt := range withOptions {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RequestTokenizerSingle(tt.name, tt.options)
			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}

	t.Run("case insensitive", func(t *testing.T) {
		got, err := RequestTokenizerSingle("WhiteSpace", nil)
		assert.NoError(t, err)
		assert.Equal(t, []string{"A", "b"}, terms(got.Tokenize([]byte("A b"))))
	})
	t.Run("option error propagates", func(t *testing.T) {
		got, err := RequestTokenizerSingle("character", map[string]interface{}{"char": "bad"})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("unknown", func(t *testing.T) {
		got, err := RequestTokenizerSingle("nope", nil)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("empty", func(t *testing.T) {
		got, err := RequestTokenizerSingle("", nil)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
