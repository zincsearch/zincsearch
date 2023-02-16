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

package api

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zinclabs/zincsearch/pkg/zutils/json"
)

func TestDocumentBulk(t *testing.T) {

	t.Run("POST /api/_bulk", func(t *testing.T) {
		t.Run("bulk documents", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(bulkData)
			resp := request("POST", "/api/_bulk", body)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
		t.Run("bulk documents with delete", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(bulkDataWithDelete)
			resp := request("POST", "/api/_bulk", body)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
		t.Run("bulk with error input", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"index":{}}`)
			resp := request("POST", "/api/_bulk", body)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
	})

	t.Run("POST /api/:target/_bulk", func(t *testing.T) {
		t.Run("bulk create documents with not exist indexName", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			data := strings.ReplaceAll(bulkData, `"_index": "games3"`, `"_index": ""`)
			body.WriteString(data)
			resp := request("POST", "/api/notExistIndex/_bulk", body)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
		t.Run("bulk create documents with exist indexName", func(t *testing.T) {
			// create index
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"name": "` + indexName + `", "storage_type": "disk"}`)
			resp := request("PUT", "/api/index", body)
			assert.Equal(t, http.StatusBadRequest, resp.Code)

			respData := make(map[string]string)
			err := json.Unmarshal(resp.Body.Bytes(), &respData)
			assert.NoError(t, err)
			assert.Equal(t, "index ["+indexName+"] already exists", respData["error"])

			// check bulk
			body.Reset()
			data := strings.ReplaceAll(bulkData, `"_index": "games3"`, `"_index": ""`)
			body.WriteString(data)
			resp = request("POST", "/api/"+indexName+"/_bulk", body)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
		t.Run("bulk with error input", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"index":{}}`)
			resp := request("POST", "/api/"+indexName+"/_bulk", body)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
	})

}
