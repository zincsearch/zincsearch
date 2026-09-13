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

package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/pkg/meta"
)

func TestRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   []interface{}
		want    []*meta.Field
		wantErr bool
	}{
		{"nil", nil, nil, false},
		{"strings", []interface{}{"a", "b"}, []*meta.Field{{Field: "a"}, {Field: "b"}}, false},
		{"object with format", []interface{}{map[string]interface{}{"Field": "ts", "FORMAT": "2006-01-02"}},
			[]*meta.Field{{Field: "ts", Format: "2006-01-02"}}, false},
		{"object ignores unknown keys", []interface{}{map[string]interface{}{"field": "ts", "other": 1}},
			[]*meta.Field{{Field: "ts"}}, false},
		{"mixed", []interface{}{"a", map[string]interface{}{"field": "b"}},
			[]*meta.Field{{Field: "a"}, {Field: "b"}}, false},
		{"unsupported type", []interface{}{1}, nil, true},
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
}

func TestResponse(t *testing.T) {
	mappings := meta.NewMappings()
	mappings.SetProperty("name", meta.NewProperty("keyword"))
	dateProp := meta.NewProperty("date")
	dateProp.Format = "2006-01-02T15:04:05Z07:00"
	mappings.SetProperty("created", dateProp)
	mappings.SetProperty("meta.updated", dateProp)

	data := []byte(`{"name":"zinc","created":"2022-03-04T05:06:07Z","meta.updated":"2022-03-04T05:06:07Z","meta.tag":"x"}`)

	tests := []struct {
		name   string
		fields []*meta.Field
		data   []byte
		want   map[string]interface{}
	}{
		{"no fields returns nil", nil, data, nil},
		{"invalid json returns nil", []*meta.Field{{Field: "name"}}, []byte(`{`), nil},
		{"plain field", []*meta.Field{{Field: "name"}}, data, map[string]interface{}{"name": []interface{}{"zinc"}}},
		{"missing field skipped", []*meta.Field{{Field: "nope"}}, data, map[string]interface{}{}},
		{"date without format is raw", []*meta.Field{{Field: "created"}}, data,
			map[string]interface{}{"created": []interface{}{"2022-03-04T05:06:07Z"}}},
		{"date with format is reformatted", []*meta.Field{{Field: "created", Format: "2006-01-02"}}, data,
			map[string]interface{}{"created": []interface{}{"2022-03-04"}}},
		{"date with unparsable value falls back to raw", []*meta.Field{{Field: "created", Format: "2006"}},
			[]byte(`{"created":"not-a-date"}`), map[string]interface{}{"created": []interface{}{"not-a-date"}}},
		{"wildcard with format", []*meta.Field{{Field: "meta.*", Format: "2006"}}, data,
			map[string]interface{}{"meta.updated": []interface{}{"2022"}, "meta.tag": []interface{}{"x"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Response(tt.fields, tt.data, mappings))
		})
	}
}
