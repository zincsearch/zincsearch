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

package index

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/pkg/core"
	"github.com/zincsearch/zincsearch/pkg/zutils/json"
	"github.com/zincsearch/zincsearch/test/utils"
)

func TestCreate(t *testing.T) {
	type args struct {
		code       int
		data       string
		params     map[string]string
		result     string
		mappingRes string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "create by json",
			args: args{
				code:   http.StatusOK,
				data:   `{"name":"TestIndexCreate.index_1","disk":"disk"}`,
				params: map[string]string{"target": ""},
				result: `"message":"ok"`,
			},
			wantErr: false,
		},
		{
			name: "create by target",
			args: args{
				code:   http.StatusOK,
				data:   `{"name":"","disk":"disk"}`,
				params: map[string]string{"target": "TestIndexCreate.index_2"},
				result: `"message":"ok"`,
			},
			wantErr: false,
		},
		{
			name: "create by error json",
			args: args{
				code:   http.StatusBadRequest,
				data:   `{"name":"TestIndexCreate.index_3"x,"disk":"disk"}`,
				params: map[string]string{"target": ""},
				result: `"error":`,
			},
			wantErr: true,
		},
		{
			name: "create by empty",
			args: args{
				code:   http.StatusBadRequest,
				data:   `{"name":"","disk":"disk"}`,
				params: map[string]string{"target": ""},
				result: "should be not empty",
			},
			wantErr: true,
		},
		{
			name: "create with analyzer",
			args: args{
				code:   http.StatusOK,
				data:   `{"name":"TestIndexCreate.index_5","disk":"disk","settings":{"analysis":{"analyzer":{"test_analyzer":{"type":"custom","tokenizer":"standard","filter":["lowercase"]}}}}}`,
				params: map[string]string{"target": ""},
				result: `"message":"ok"`,
			},
			wantErr: false,
		},
		{
			name: "create with sub-fields",
			args: args{
				code:       http.StatusOK,
				data:       `{"name":"TestIndexCreate.index_6","disk":"disk","mappings":{"properties":{"@timestamp":{"type":"date"},"Athlete":{"type":"text","fields":{"my_keyword":{"type":"keyword"}}}}}}`,
				mappingRes: `"Athlete.my_keyword":{"type":"keyword"`,
				params:     map[string]string{"target": ""},
				result:     `"message":"ok"`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			utils.SetGinRequestData(c, tt.args.data)
			utils.SetGinRequestParams(c, tt.args.params)
			Create(c)
			assert.Equal(t, tt.args.code, w.Code)
			assert.Contains(t, w.Body.String(), tt.args.result)

			resp := make(map[string]string)
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tt.args.mappingRes != "" {
				c, w := utils.NewGinContext()
				utils.SetGinRequestParams(c, map[string]string{"target": resp["index"]})

				GetMapping(c)
				assert.Equal(t, tt.args.code, w.Code)
				assert.Contains(t, w.Body.String(), tt.args.mappingRes)
			}

			if !tt.wantErr {
				err = core.DeleteIndex(resp["index"])
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateES(t *testing.T) {
	type args struct {
		code   int
		data   string
		params map[string]string
		result string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "create by json",
			args: args{
				code:   http.StatusOK,
				data:   `{"name":"TestIndexCreate.index_1"}`,
				params: map[string]string{"target": ""},
				result: `"acknowledged":true`,
			},
			wantErr: false,
		},
		{
			name: "create by target",
			args: args{
				code:   http.StatusOK,
				data:   `{"name":""}`,
				params: map[string]string{"target": "TestIndexCreate.index_2"},
				result: `"acknowledged":true`,
			},
			wantErr: false,
		},
		{
			name: "create by target",
			args: args{
				code:   http.StatusOK,
				data:   `{"name":""}`,
				params: map[string]string{"target": "TestIndexCreate.index_3"},
				result: `"acknowledged":true`,
			},
			wantErr: false,
		},
		{
			name: "create by error json",
			args: args{
				code:   http.StatusBadRequest,
				data:   `{"name":"TestIndexCreate.index_4"x}`,
				params: map[string]string{"target": ""},
				result: `"error":`,
			},
			wantErr: true,
		},
		{
			name: "create by empty",
			args: args{
				code:   http.StatusBadRequest,
				data:   `{"name":"","disk":"disk"}`,
				params: map[string]string{"target": ""},
				result: "should be not empty",
			},
			wantErr: true,
		},
		{
			name: "create with analyzer",
			args: args{
				code:   http.StatusOK,
				data:   `{"name":"TestIndexCreate.index_6","disk":"disk","settings":{"analysis":{"analyzer":{"test_analyzer":{"type":"custom","tokenizer":"standard","filter":["lowercase"]}}}}}`,
				params: map[string]string{"target": ""},
				result: `"acknowledged":true`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			utils.SetGinRequestData(c, tt.args.data)
			utils.SetGinRequestParams(c, tt.args.params)
			CreateES(c)
			assert.Equal(t, tt.args.code, w.Code)
			assert.Contains(t, w.Body.String(), tt.args.result)

			resp := make(map[string]interface{})
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if !tt.wantErr {
				err = core.DeleteIndex(resp["index"].(string))
				assert.NoError(t, err)
			}
		})
	}
}
