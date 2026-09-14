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

package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/pkg/uquery/analysis/internal/testutil"
)

func TestNewCharGroupTokenizer(t *testing.T) {
	t.Run("split on chars", func(t *testing.T) {
		z, err := NewCharGroupTokenizer(map[string]interface{}{"tokenize_on_chars": []interface{}{"-", "whitespace"}})
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, testutil.Terms(z.Tokenize([]byte("a-b c"))))
	})
	t.Run("missing option", func(t *testing.T) {
		z, err := NewCharGroupTokenizer(nil)
		assert.Error(t, err)
		assert.Nil(t, z)
	})
}

func TestNewCharacterTokenizer(t *testing.T) {
	tests := []struct {
		name    string
		char    string
		text    string
		want    []string
		wantErr bool
	}{
		{name: "letter", char: "letter", text: "ab1cd 2ef", want: []string{"ab", "cd", "ef"}},
		{name: "number", char: "number", text: "ab1cd 22ef", want: []string{"1", "22"}},
		{name: "digit alias", char: "digit", text: "a1b2", want: []string{"1", "2"}},
		{name: "punct", char: "punct", text: "a,b.c", want: []string{",", "."}},
		{name: "punctuation alias", char: "punctuation", text: "a,b", want: []string{","}},
		{name: "space", char: "space", text: "a b", want: []string{" "}},
		{name: "whitespace alias", char: "whitespace", text: "a b", want: []string{" "}},
		{name: "white_space alias", char: "white_space", text: "a b", want: []string{" "}},
		{name: "symbol", char: "symbol", text: "a+b", want: []string{"+"}},
		{name: "graphic", char: "graphic", text: "ab\ncd", want: []string{"ab", "cd"}},
		{name: "print", char: "print", text: "ab\ncd", want: []string{"ab", "cd"}},
		{name: "control", char: "control", text: "ab\ncd", want: []string{"\n"}},
		{name: "mark", char: "mark", text: "e\u0301a", want: []string{"\u0301"}},
		{name: "unknown", char: "emoji", wantErr: true},
		{name: "missing", char: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z, err := NewCharacterTokenizer(map[string]interface{}{"char": tt.char})
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, z)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, testutil.Terms(z.Tokenize([]byte(tt.text))))
		})
	}
}

func TestNewEdgeNgramTokenizer(t *testing.T) {
	tests := []struct {
		name    string
		options interface{}
		want    []string
		wantErr bool
	}{
		{name: "default 1-2", options: nil, want: []string{"a", "ab"}},
		{name: "2-3", options: map[string]interface{}{"min_gram": float64(2), "max_gram": float64(3)}, want: []string{"ab", "abc"}},
		{name: "min greater than max", options: map[string]interface{}{"min_gram": float64(3), "max_gram": float64(2)}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z, err := NewEdgeNgramTokenizer(tt.options)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, z)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, testutil.Terms(z.Tokenize([]byte("abcd"))))
		})
	}
}

func TestNewExceptionTokenizer(t *testing.T) {
	t.Run("keeps matched patterns whole", func(t *testing.T) {
		z, err := NewExceptionTokenizer(map[string]interface{}{"patterns": []interface{}{`[a-z]+\.[a-z]+`}})
		assert.NoError(t, err)
		assert.Equal(t, []string{"visit", "zinc.com", "now"}, testutil.Terms(z.Tokenize([]byte("visit zinc.com now"))))
	})
	t.Run("missing patterns", func(t *testing.T) {
		z, err := NewExceptionTokenizer(nil)
		assert.Error(t, err)
		assert.Nil(t, z)
	})
	t.Run("invalid pattern", func(t *testing.T) {
		z, err := NewExceptionTokenizer(map[string]interface{}{"patterns": []interface{}{"("}})
		assert.Error(t, err)
		assert.Nil(t, z)
	})
}

func TestNewLowerCaseTokenizer(t *testing.T) {
	z, err := NewLowerCaseTokenizer()
	assert.NoError(t, err)
	assert.Equal(t, []string{"hello", "world"}, testutil.Terms(z.Tokenize([]byte("Hello World"))))
}

func TestNewNgramTokenizer(t *testing.T) {
	tests := []struct {
		name    string
		options interface{}
		want    []string
		wantErr bool
	}{
		{name: "default 1-2", options: nil, want: []string{"a", "ab", "b", "bc", "c"}},
		{name: "3-3", options: map[string]interface{}{"min_gram": float64(3), "max_gram": float64(3)}, want: []string{"abc"}},
		{name: "min greater than max", options: map[string]interface{}{"min_gram": float64(3), "max_gram": float64(2)}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z, err := NewNgramTokenizer(tt.options)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, z)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, testutil.Terms(z.Tokenize([]byte("abc"))))
		})
	}
}

func TestNewPathHierarchyTokenizer(t *testing.T) {
	tests := []struct {
		name    string
		options interface{}
		text    string
		want    []string
	}{
		{name: "default", options: nil, text: "/a/b/c", want: []string{"/a", "/a/b", "/a/b/c"}},
		{
			name:    "custom delimiter and replacement",
			options: map[string]interface{}{"delimiter": "-", "replacement": "/"},
			text:    "a-b-c",
			want:    []string{"a", "a/b", "a/b/c"},
		},
		{
			name:    "skip",
			options: map[string]interface{}{"skip": float64(1)},
			text:    "/a/b/c",
			want:    []string{"/b", "/b/c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z, err := NewPathHierarchyTokenizer(tt.options)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, testutil.Terms(z.Tokenize([]byte(tt.text))))
		})
	}
}

func TestNewRegexpTokenizer(t *testing.T) {
	tests := []struct {
		name    string
		options interface{}
		want    []string
		wantErr bool
	}{
		{name: "default word pattern", options: nil, want: []string{"a1", "b", "c"}},
		{name: "custom", options: map[string]interface{}{"pattern": "[a-z]+"}, want: []string{"a", "b", "c"}},
		{name: "invalid", options: map[string]interface{}{"pattern": "("}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z, err := NewRegexpTokenizer(tt.options)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, z)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, testutil.Terms(z.Tokenize([]byte("a1-b c"))))
		})
	}
}
