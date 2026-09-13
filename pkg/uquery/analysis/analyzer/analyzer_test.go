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

package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vcaesar/riot/analysis"
)

func terms(ts analysis.TokenStream) []string {
	out := make([]string, 0, len(ts))
	for _, t := range ts {
		out = append(out, string(t.Term))
	}
	return out
}

func TestNewRegexpAnalyzer(t *testing.T) {
	tests := []struct {
		name    string
		options interface{}
		text    string
		want    []string
		wantErr bool
	}{
		{
			name:    "default pattern lowercases",
			options: nil,
			text:    "Hello World-Foo",
			want:    []string{"hello", "world", "foo"},
		},
		{
			name:    "custom pattern keep case",
			options: map[string]interface{}{"pattern": "[A-Z]+", "lowercase": false},
			text:    "Hello World",
			want:    []string{"H", "W"},
		},
		{
			name:    "stopwords",
			options: map[string]interface{}{"stopwords": []interface{}{"world"}},
			text:    "Hello World",
			want:    []string{"hello"},
		},
		{
			name:    "invalid pattern",
			options: map[string]interface{}{"pattern": "("},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ana, err := NewRegexpAnalyzer(tt.options)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, ana)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, terms(ana.Analyze([]byte(tt.text))))
		})
	}
}

func TestNewStandardAnalyzer(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		ana, err := NewStandardAnalyzer(nil)
		assert.NoError(t, err)
		assert.Equal(t, []string{"hello", "world"}, terms(ana.Analyze([]byte("Hello World"))))
	})
	t.Run("stopwords", func(t *testing.T) {
		ana, err := NewStandardAnalyzer(map[string]interface{}{"stopwords": []string{"hello"}})
		assert.NoError(t, err)
		assert.Equal(t, []string{"world"}, terms(ana.Analyze([]byte("Hello World"))))
	})
}

func TestNewStopAnalyzer(t *testing.T) {
	t.Run("default english stopwords", func(t *testing.T) {
		ana, err := NewStopAnalyzer(nil)
		assert.NoError(t, err)
		assert.Equal(t, []string{"quick", "fox"}, terms(ana.Analyze([]byte("The quick fox"))))
	})
	t.Run("custom stopwords", func(t *testing.T) {
		ana, err := NewStopAnalyzer(map[string]interface{}{"stopwords": []interface{}{"quick"}})
		assert.NoError(t, err)
		assert.Equal(t, []string{"the", "fox"}, terms(ana.Analyze([]byte("The quick fox"))))
	})
}

func TestNewWhitespaceAnalyzer(t *testing.T) {
	ana, err := NewWhitespaceAnalyzer()
	assert.NoError(t, err)
	assert.Equal(t, []string{"Hello", "World-Foo"}, terms(ana.Analyze([]byte("Hello  World-Foo"))))
}
