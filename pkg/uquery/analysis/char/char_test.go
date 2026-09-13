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

package char

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMappingCharFilter(t *testing.T) {
	tests := []struct {
		name    string
		options interface{}
		text    string
		want    string
		wantErr bool
	}{
		{
			name:    "replace",
			options: map[string]interface{}{"mappings": []interface{}{"a => 1", "b => 2"}},
			text:    "abc",
			want:    "12c",
		},
		{
			name:    "nil options",
			options: nil,
			wantErr: true,
		},
		{
			name:    "empty mappings",
			options: map[string]interface{}{"mappings": []interface{}{}},
			wantErr: true,
		},
		{
			name:    "invalid rule",
			options: map[string]interface{}{"mappings": []interface{}{"a=>1"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewMappingCharFilter(tt.options)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, f)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(f.Filter([]byte(tt.text))))
		})
	}
}

func TestNewRegexpCharFilter(t *testing.T) {
	tests := []struct {
		name    string
		options interface{}
		text    string
		want    string
		wantErr bool
	}{
		{
			name:    "replace digits",
			options: map[string]interface{}{"pattern": "[0-9]+", "replacement": "#"},
			text:    "a1b22c",
			want:    "a#b#c",
		},
		{
			name:    "no replacement removes match",
			options: map[string]interface{}{"pattern": "-"},
			text:    "a-b",
			want:    "ab",
		},
		{
			name:    "missing pattern",
			options: map[string]interface{}{},
			wantErr: true,
		},
		{
			name:    "empty pattern",
			options: map[string]interface{}{"pattern": ""},
			wantErr: true,
		},
		{
			name:    "invalid pattern",
			options: map[string]interface{}{"pattern": "["},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewRegexpCharFilter(tt.options)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, f)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(f.Filter([]byte(tt.text))))
		})
	}
}

func TestNewSTConvertCharFilter(t *testing.T) {
	t.Run("s2t", func(t *testing.T) {
		f, err := NewSTConvertCharFilter(map[string]interface{}{"convert_type": "s2t"})
		assert.NoError(t, err)
		assert.Equal(t, "漢字", string(f.Filter([]byte("汉字"))))
	})
	t.Run("missing convert_type", func(t *testing.T) {
		f, err := NewSTConvertCharFilter(nil)
		assert.Error(t, err)
		assert.Nil(t, f)
	})
	t.Run("invalid convert_type", func(t *testing.T) {
		f, err := NewSTConvertCharFilter(map[string]interface{}{"convert_type": "xx"})
		assert.Error(t, err)
		assert.Nil(t, f)
	})
}
