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

	"github.com/zincsearch/zincsearch/test/utils"
)

func TestTemplate(t *testing.T) {
	t.Run("create template", func(t *testing.T) {
		type args struct {
			code    int
			data    map[string]interface{}
			rawData string
			target  string
			result  string
		}
		tests := []struct {
			name    string
			args    args
			wantErr bool
		}{
			{
				name: "normal",
				args: args{
					code: http.StatusOK,
					data: map[string]interface{}{
						"name":           "test1",
						"index_patterns": []string{"log-*"},
						"priority":       1,
						"template": map[string]interface{}{
							"settings": map[string]interface{}{
								"number_of_shards":   3,
								"number_of_replicas": 1,
							},
							"mappings": map[string]interface{}{
								"properties": map[string]interface{}{
									"Athlete": map[string]interface{}{
										"type":          "text",
										"index":         true,
										"store":         false,
										"sortable":      false,
										"aggregatable":  false,
										"highlightable": false,
									},
								},
							},
						},
					},
					target: "TestTemplate.index_1",
					result: `{"message":"ok"`,
				},
				wantErr: false,
			},
			{
				name: "empty",
				args: args{
					code:    http.StatusBadRequest,
					rawData: `{}`,
					target:  "",
					result:  `should be not empty`,
				},
				wantErr: false,
			},
			{
				name: "with err json",
				args: args{
					code:    http.StatusBadRequest,
					rawData: `{"x":x}`,
					target:  "TestTemplate.index_3",
					result:  `invalid character`,
				},
				wantErr: false,
			},
			{
				name: "with err type",
				args: args{
					code: http.StatusBadRequest,
					data: map[string]interface{}{
						"name":           "test2",
						"index_patterns": "log-*",
						"priority":       1,
						"template": map[string]interface{}{
							"settings": map[string]interface{}{
								"number_of_shards":   3,
								"number_of_replicas": 1,
							},
						},
					},
					target: "TestTemplate.index_4",
					result: `index_patterns value should be an array of string`,
				},
				wantErr: false,
			},
			{
				name: "with same priority",
				args: args{
					code: http.StatusBadRequest,
					data: map[string]interface{}{
						"name":           "test3",
						"index_patterns": []string{"log-*"},
						"priority":       1,
						"template": map[string]interface{}{
							"settings": map[string]interface{}{
								"number_of_shards":   3,
								"number_of_replicas": 1,
							},
						},
					},
					target: "TestTemplate.index_4",
					result: `have the same priority`,
				},
				wantErr: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c, w := utils.NewGinContext()
				if tt.args.data != nil {
					utils.SetGinRequestData(c, tt.args.data)
				}
				if tt.args.rawData != "" {
					utils.SetGinRequestData(c, tt.args.rawData)
				}
				utils.SetGinRequestParams(c, map[string]string{"target": tt.args.target})
				CreateTemplate(c)
				assert.Equal(t, tt.args.code, w.Code)
				assert.Contains(t, w.Body.String(), tt.args.result)
			})
		}
	})

	t.Run("get template", func(t *testing.T) {
		type args struct {
			code   int
			target string
			result string
		}
		tests := []struct {
			name    string
			args    args
			wantErr bool
		}{
			{
				name: "normal",
				args: args{
					code:   http.StatusOK,
					target: "TestTemplate.index_1",
					result: `"index_patterns":`,
				},
				wantErr: false,
			},
			{
				name: "empty",
				args: args{
					code:   http.StatusBadRequest,
					target: "",
					result: `should be not empty`,
				},
				wantErr: false,
			},
			{
				name: "not exists",
				args: args{
					code:   http.StatusNotFound,
					target: "TestTemplate.index_N",
					result: `does not exists`,
				},
				wantErr: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c, w := utils.NewGinContext()
				utils.SetGinRequestParams(c, map[string]string{"target": tt.args.target})
				GetTemplate(c)
				assert.Equal(t, tt.args.code, w.Code)
				assert.Contains(t, w.Body.String(), tt.args.result)
			})
		}
	})

	t.Run("list template", func(t *testing.T) {
		type args struct {
			code   int
			target string
			result string
		}
		tests := []struct {
			name    string
			args    args
			wantErr bool
		}{
			{
				name: "pattern",
				args: args{
					code:   http.StatusOK,
					target: "log-*",
					result: `"index_patterns":`,
				},
				wantErr: false,
			},
			{
				name: "empty",
				args: args{
					code:   http.StatusOK,
					target: "",
					result: `"index_patterns":`,
				},
				wantErr: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c, w := utils.NewGinContext()
				utils.SetGinRequestURL(c, "", map[string]string{"pattern": tt.args.target})
				ListTemplate(c)
				assert.Equal(t, tt.args.code, w.Code)
				assert.Contains(t, w.Body.String(), tt.args.result)
			})
		}
	})

	t.Run("delete template", func(t *testing.T) {
		type args struct {
			code   int
			target string
			result string
		}
		tests := []struct {
			name    string
			args    args
			wantErr bool
		}{
			{
				name: "pattern",
				args: args{
					code:   http.StatusOK,
					target: "test1",
					result: `"message":"ok"`,
				},
				wantErr: false,
			},
			{
				name: "empty",
				args: args{
					code:   http.StatusOK,
					target: "",
					result: `"message":"ok"`,
				},
				wantErr: false,
			},
			{
				name: "not exists",
				args: args{
					code:   http.StatusOK,
					target: "testN",
					result: `"message":"ok"`,
				},
				wantErr: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c, w := utils.NewGinContext()
				utils.SetGinRequestParams(c, map[string]string{"target": tt.args.target})
				DeleteTemplate(c)
				assert.Equal(t, tt.args.code, w.Code)
				assert.Contains(t, w.Body.String(), tt.args.result)
			})
		}
	})
}
