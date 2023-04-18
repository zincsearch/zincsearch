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

package document

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/pkg/core"
	"github.com/zincsearch/zincsearch/test/utils"
)

func TestCreateUpdate(t *testing.T) {
	type args struct {
		code    int
		data    map[string]interface{}
		rawData string
		params  map[string]string
		result  string
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
					"_id":  "1",
					"name": "user",
					"role": "create",
				},
				params: map[string]string{
					"target": "TestDocumentCreateUpdate.index_1",
				},
				result: `"id":"1"`,
			},
		},
		{
			name: "update 1",
			args: args{
				code: http.StatusOK,
				data: map[string]interface{}{
					"_id":  "1",
					"name": "user",
					"role": "update1",
				},
				params: map[string]string{
					"target": "TestDocumentCreateUpdate.index_1",
				},
				result: `"id":"1"`,
			},
		},
		{
			name: "update 2",
			args: args{
				code: http.StatusOK,
				data: map[string]interface{}{
					"name": "user",
					"role": "update2",
				},
				params: map[string]string{
					"target": "TestDocumentCreateUpdate.index_1",
					"id":     "1",
				},
				result: `"id":"1"`,
			},
		},
		{
			name: "generate id",
			args: args{
				code: http.StatusOK,
				data: map[string]interface{}{
					"name": "user",
					"role": "update2",
				},
				params: map[string]string{
					"target": "TestDocumentCreateUpdate.index_1",
				},
				result: `"id":`,
			},
		},
		{
			name: "error json",
			args: args{
				code:    http.StatusBadRequest,
				rawData: `{"name": "user", "role": "create", }`,
				params: map[string]string{
					"target": "TestDocumentCreateUpdate.index_1",
				},
				result: `"error":`,
			},
		},
		{
			name: "empty index",
			args: args{
				code:    http.StatusInternalServerError,
				rawData: `{"name": "user", "role": "create"}`,
				params:  map[string]string{},
				result:  "index name cannot be empty",
			},
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
			if tt.args.params != nil {
				utils.SetGinRequestParams(c, tt.args.params)
			}
			CreateUpdate(c)
			assert.Equal(t, tt.args.code, w.Code)
			assert.Contains(t, w.Body.String(), tt.args.result)
		})
	}

	t.Run("cleanup", func(t *testing.T) {
		err := core.DeleteIndex("TestDocumentCreateUpdate.index_1")
		assert.NoError(t, err)
	})
}
