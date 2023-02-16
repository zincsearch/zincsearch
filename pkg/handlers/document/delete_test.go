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
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zincsearch/pkg/core"
	"github.com/zinclabs/zincsearch/test/utils"
)

func TestDelete(t *testing.T) {
	type args struct {
		code   int
		params map[string]string
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
				params: map[string]string{
					"target": "TestDocumentDelete.index_1",
					"id":     "1",
				},
				result: `"id":"1"`,
			},
		},
		{
			name: "empty id",
			args: args{
				code: http.StatusBadRequest,
				params: map[string]string{
					"target": "TestDocumentDelete.index_1",
				},
				result: `"id is empty"`,
			},
		},
		{
			name: "not exists id",
			args: args{
				code: http.StatusBadRequest,
				params: map[string]string{
					"target": "TestDocumentDelete.index_1",
					"id":     "2",
				},
				result: `"id not found"`,
			},
		},
		{
			name: "not exists index",
			args: args{
				code: http.StatusBadRequest,
				params: map[string]string{
					"target": "TestDocumentDelete.index_2",
					"id":     "1",
				},
				result: "index does not exists",
			},
		},
	}

	// create a document
	indexName := "TestDocumentDelete.index_1"
	t.Run("prepare", func(t *testing.T) {
		data := map[string]interface{}{
			"_id":  "1",
			"name": "user",
			"role": "create",
		}
		params := map[string]string{
			"target": indexName,
		}

		c, w := utils.NewGinContext()
		utils.SetGinRequestData(c, data)
		utils.SetGinRequestParams(c, params)
		CreateUpdate(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"id":"1"`)

		// wait for WAL write to index
		time.Sleep(time.Second)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			utils.SetGinRequestParams(c, tt.args.params)
			Delete(c)
			assert.Equal(t, tt.args.code, w.Code)
			assert.Contains(t, w.Body.String(), tt.args.result)
		})
	}

	t.Run("cleanup", func(t *testing.T) {
		err := core.DeleteIndex(indexName)
		assert.NoError(t, err)
	})
}
