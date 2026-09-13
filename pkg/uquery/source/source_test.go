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

package source

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/pkg/meta"
)

func TestRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    *meta.Source
		wantErr bool
	}{
		{"nil enables source", nil, &meta.Source{Enable: true}, false},
		{"bool false", false, &meta.Source{Enable: false}, false},
		{"bool true", true, &meta.Source{Enable: true}, false},
		{"field list", []interface{}{"name", "age"}, &meta.Source{Enable: true, Fields: []string{"name", "age"}}, false},
		{"empty list", []interface{}{}, &meta.Source{Enable: true, Fields: []string{}}, false},
		{"non-string in list", []interface{}{"name", 1}, nil, true},
		{"unsupported type", "name", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Request(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}

	t.Run("passthrough meta.Source", func(t *testing.T) {
		in := &meta.Source{Enable: false, Fields: []string{"x"}}
		got, err := Request(in)
		assert.NoError(t, err)
		assert.Same(t, in, got)
	})
}

func TestResponse(t *testing.T) {
	data := []byte(`{"name":"zinc","age":3,"addr.city":"sf","addr.zip":"94103"}`)
	tests := []struct {
		name   string
		source *meta.Source
		data   []byte
		want   map[string]interface{}
	}{
		{"disabled returns empty", &meta.Source{Enable: false}, data, map[string]interface{}{}},
		{"invalid json returns nil", &meta.Source{Enable: true}, []byte(`{`), nil},
		{"no fields returns all", &meta.Source{Enable: true}, data,
			map[string]interface{}{"name": "zinc", "age": float64(3), "addr.city": "sf", "addr.zip": "94103"}},
		{"exact fields", &meta.Source{Enable: true, Fields: []string{"name", "missing"}}, data,
			map[string]interface{}{"name": "zinc"}},
		{"wildcard fields", &meta.Source{Enable: true, Fields: []string{"addr.*"}}, data,
			map[string]interface{}{"addr.city": "sf", "addr.zip": "94103"}},
		{"mixed exact and wildcard", &meta.Source{Enable: true, Fields: []string{"age", "addr.*"}}, data,
			map[string]interface{}{"age": float64(3), "addr.city": "sf", "addr.zip": "94103"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Response(tt.source, tt.data))
		})
	}
}
