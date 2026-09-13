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

package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vcaesar/riot/analysis"
	"github.com/vcaesar/riot/analysis/tokenizer"
)

func tokens(text string) analysis.TokenStream {
	return tokenizer.NewWhitespaceTokenizer().Tokenize([]byte(text))
}

func terms(ts analysis.TokenStream) []string {
	out := make([]string, 0, len(ts))
	for _, t := range ts {
		out = append(out, string(t.Term))
	}
	return out
}

func TestNewDictTokenFilter(t *testing.T) {
	t.Run("decompound", func(t *testing.T) {
		f, err := NewDictTokenFilter(map[string]interface{}{"words": []interface{}{"soft", "ball"}})
		assert.NoError(t, err)
		assert.Equal(t, []string{"softball", "soft", "ball"}, terms(f.Filter(tokens("softball"))))
	})
	t.Run("min_word_size skips short words", func(t *testing.T) {
		f, err := NewDictTokenFilter(map[string]interface{}{
			"words":         []interface{}{"soft", "ball"},
			"min_word_size": float64(20),
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"softball"}, terms(f.Filter(tokens("softball"))))
	})
	t.Run("missing words", func(t *testing.T) {
		f, err := NewDictTokenFilter(nil)
		assert.Error(t, err)
		assert.Nil(t, f)
	})
	t.Run("words not strings", func(t *testing.T) {
		f, err := NewDictTokenFilter(map[string]interface{}{"words": []interface{}{1}})
		assert.Error(t, err)
		assert.Nil(t, f)
	})
}

func TestNewEdgeNgramTokenFilter(t *testing.T) {
	tests := []struct {
		name    string
		options interface{}
		want    []string
		wantErr bool
	}{
		{name: "default front 1-2", options: nil, want: []string{"a", "ab"}},
		{name: "front 2-3", options: map[string]interface{}{"min_gram": float64(2), "max_gram": float64(3)}, want: []string{"ab", "abc"}},
		{name: "back", options: map[string]interface{}{"side": "BACK", "min_gram": float64(2), "max_gram": float64(2)}, want: []string{"cd"}},
		{name: "min greater than max", options: map[string]interface{}{"min_gram": float64(3), "max_gram": float64(2)}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewEdgeNgramTokenFilter(tt.options)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, f)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, terms(f.Filter(tokens("abcd"))))
		})
	}
}

func TestNewElisionTokenFilter(t *testing.T) {
	t.Run("default french articles", func(t *testing.T) {
		f, err := NewElisionTokenFilter(nil)
		assert.NoError(t, err)
		assert.Equal(t, []string{"avion", "homme"}, terms(f.Filter(tokens("l'avion d'homme"))))
	})
	t.Run("custom articles", func(t *testing.T) {
		f, err := NewElisionTokenFilter(map[string]interface{}{"articles": []interface{}{"x"}})
		assert.NoError(t, err)
		assert.Equal(t, []string{"avion", "l'homme"}, terms(f.Filter(tokens("x'avion l'homme"))))
	})
}

func TestNewKeywordTokenFilter(t *testing.T) {
	t.Run("marks keywords", func(t *testing.T) {
		f, err := NewKeywordTokenFilter(map[string]interface{}{"keywords": []interface{}{"foo"}})
		assert.NoError(t, err)
		out := f.Filter(tokens("foo bar"))
		assert.True(t, out[0].KeyWord)
		assert.False(t, out[1].KeyWord)
	})
	t.Run("missing keywords", func(t *testing.T) {
		f, err := NewKeywordTokenFilter(nil)
		assert.Error(t, err)
		assert.Nil(t, f)
	})
}

func TestNewLengthTokenFilter(t *testing.T) {
	t.Run("default 1-2", func(t *testing.T) {
		f, err := NewLengthTokenFilter(nil)
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "bb"}, terms(f.Filter(tokens("a bb ccc"))))
	})
	t.Run("custom", func(t *testing.T) {
		f, err := NewLengthTokenFilter(map[string]interface{}{"min": float64(2), "max": float64(3)})
		assert.NoError(t, err)
		assert.Equal(t, []string{"bb", "ccc"}, terms(f.Filter(tokens("a bb ccc dddd"))))
	})
}

func TestNewNgramTokenFilter(t *testing.T) {
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
			f, err := NewNgramTokenFilter(tt.options)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, f)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, terms(f.Filter(tokens("abc"))))
		})
	}
}

func TestNewRegexpTokenFilter(t *testing.T) {
	tests := []struct {
		name    string
		options interface{}
		want    []string
		wantErr bool
	}{
		{name: "replace", options: map[string]interface{}{"pattern": "o+", "replacement": "0"}, want: []string{"f0", "bar"}},
		{name: "remove", options: map[string]interface{}{"pattern": "o"}, want: []string{"f", "bar"}},
		{name: "missing pattern", options: nil, wantErr: true},
		{name: "empty pattern", options: map[string]interface{}{"pattern": ""}, wantErr: true},
		{name: "invalid pattern", options: map[string]interface{}{"pattern": "("}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewRegexpTokenFilter(tt.options)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, f)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, terms(f.Filter(tokens("foo bar"))))
		})
	}
}

func TestNewShingleTokenFilter(t *testing.T) {
	tests := []struct {
		name    string
		options interface{}
		want    []string
	}{
		{name: "default", options: nil, want: []string{"a", "b", "a b", "c", "b c"}},
		{
			name:    "custom separator no original",
			options: map[string]interface{}{"token_separator": "_", "output_original": false},
			want:    []string{"a_b", "b_c"},
		},
		{
			name:    "size 3",
			options: map[string]interface{}{"min_shingle_size": float64(3), "max_shingle_size": float64(3), "output_original": false},
			want:    []string{"a b c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewShingleTokenFilter(tt.options)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, terms(f.Filter(tokens("a b c"))))
		})
	}
}

func TestNewStopTokenFilter(t *testing.T) {
	t.Run("default english", func(t *testing.T) {
		f, err := NewStopTokenFilter(nil)
		assert.NoError(t, err)
		assert.Equal(t, []string{"quick", "fox"}, terms(f.Filter(tokens("the quick fox"))))
	})
	t.Run("custom", func(t *testing.T) {
		f, err := NewStopTokenFilter(map[string]interface{}{"stopwords": []interface{}{"fox"}})
		assert.NoError(t, err)
		assert.Equal(t, []string{"the", "quick"}, terms(f.Filter(tokens("the quick fox"))))
	})
}

func TestNewTrimTokenFilter(t *testing.T) {
	f, err := NewTrimTokenFilter()
	assert.NoError(t, err)
	ts := analysis.TokenStream{{Term: []byte("  foo "), Type: analysis.AlphaNumeric}}
	assert.Equal(t, []string{"foo"}, terms(f.Filter(ts)))
}

func TestNewTruncateTokenFilter(t *testing.T) {
	t.Run("default 10", func(t *testing.T) {
		f, err := NewTruncateTokenFilter(nil)
		assert.NoError(t, err)
		assert.Equal(t, []string{"abcdefghij"}, terms(f.Filter(tokens("abcdefghijklmn"))))
	})
	t.Run("custom", func(t *testing.T) {
		f, err := NewTruncateTokenFilter(map[string]interface{}{"length": float64(3)})
		assert.NoError(t, err)
		assert.Equal(t, []string{"abc", "de"}, terms(f.Filter(tokens("abcdef de"))))
	})
}

func TestNewUnicodenormTokenFilter(t *testing.T) {
	// "é" as e + combining acute accent
	decomposed := "e\u0301"
	composed := "\u00e9"

	tests := []struct {
		name    string
		form    string
		in      string
		want    string
		wantErr bool
	}{
		{name: "NFC", form: "nfc", in: decomposed, want: composed},
		{name: "NFD", form: "NFD", in: composed, want: decomposed},
		{name: "NFKC", form: "nfkc", in: "\ufb01", want: "fi"},
		{name: "NFKD", form: "NFKD", in: "\ufb01", want: "fi"},
		{name: "unknown", form: "xyz", wantErr: true},
		{name: "missing", form: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewUnicodenormTokenFilter(map[string]interface{}{"form": tt.form})
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, f)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, []string{tt.want}, terms(f.Filter(tokens(tt.in))))
		})
	}
}

func TestNewUpperCaseTokenFilter(t *testing.T) {
	f, err := NewUpperCaseTokenFilter()
	assert.NoError(t, err)
	assert.Equal(t, []string{"FOO", "BAR"}, terms(f.Filter(tokens("foo Bar"))))
}
