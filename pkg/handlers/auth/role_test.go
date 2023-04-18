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

package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/test/utils"
)

func TestCreateUpdateRole(t *testing.T) {
	type args struct {
		code   int
		data   map[string]interface{}
		result string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				code: http.StatusOK,
				data: map[string]interface{}{
					"_id":        "role",
					"name":       "role",
					"permission": []string{"test"},
				},
				result: "message",
			},
		},
		{
			name: "error",
			args: args{
				code: http.StatusBadRequest,
				data: map[string]interface{}{
					"_id": "",
				},
				result: "error",
			},
		},
		{
			name: "empty",
			args: args{
				code: http.StatusBadRequest,
				data: map[string]interface{}{
					"_id": 123,
				},
				result: "error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			utils.SetGinRequestData(c, tt.args.data)
			CreateUpdateRole(c)
			assert.Equal(t, tt.args.code, w.Code)
			assert.Contains(t, w.Body.String(), tt.args.result)
		})
	}
}

func TestListRole(t *testing.T) {
	type args struct {
		code   int
		result string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				code:   http.StatusOK,
				result: `[{"_id":`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			ListRole(c)
			assert.Equal(t, tt.args.code, w.Code)
			assert.Contains(t, w.Body.String(), tt.args.result)
		})
	}
}

func TestDeleteRole(t *testing.T) {
	type args struct {
		code   int
		data   map[string]string
		result string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				code:   http.StatusOK,
				data:   map[string]string{"id": "1"},
				result: "message",
			},
		},
		{
			name: "empty",
			args: args{
				code:   http.StatusOK,
				data:   map[string]string{"id": ""},
				result: "message",
			},
		},
		{
			name: "nil",
			args: args{
				code:   http.StatusOK,
				data:   nil,
				result: "message",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			utils.SetGinRequestParams(c, tt.args.data)
			DeleteRole(c)
			assert.Equal(t, tt.args.code, w.Code)
			assert.Contains(t, w.Body.String(), tt.args.result)
		})
	}
}
