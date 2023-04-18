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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zincsearch/zincsearch/pkg/zutils/json"
)

func TestApiES(t *testing.T) {
	t.Run("test es api", func(t *testing.T) {
		t.Run("POST /es/_bulk", func(t *testing.T) {
			t.Run("bulk documents", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(bulkData)
				resp := request("POST", "/es/_bulk", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("bulk documents with delete", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(bulkDataWithDelete)
				resp := request("POST", "/es/_bulk", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("bulk with error input", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"index":{}}`)
				resp := request("POST", "/es/_bulk", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
		})

		t.Run("POST /es/:target/_bulk", func(t *testing.T) {
			t.Run("bulk create documents with not exist indexName", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				data := strings.ReplaceAll(bulkData, `"_index": "games3"`, `"_index": ""`)
				body.WriteString(data)
				resp := request("POST", "/es/notExistIndex/_bulk", body)
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
				resp = request("POST", "/es/"+indexName+"/_bulk", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("bulk with error input", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"index":{}}`)
				resp := request("POST", "/es/"+indexName+"/_bulk", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
		})

		t.Run("POST /es/:target/_doc", func(t *testing.T) {
			_id := ""
			t.Run("create document with not exist indexName", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/notExistIndex/_doc", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("create document with exist indexName", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_doc", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("create document with exist indexName not exist id", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_doc", body)
				assert.Equal(t, http.StatusOK, resp.Code)

				data := make(map[string]interface{})
				err := json.Unmarshal(resp.Body.Bytes(), &data)
				assert.NoError(t, err)
				assert.NotEqual(t, "", data["id"])
				_id = data["id"].(string)
			})
			t.Run("update document with exist indexName and exist id", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				data := strings.Replace(indexData, "{", "{\"_id\": \""+_id+"\",", 1)
				body.WriteString(data)
				resp := request("POST", "/es/"+indexName+"/_doc", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("create document with error input", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(`data`)
				resp := request("POST", "/es/"+indexName+"/_doc", body)
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			})
		})

		t.Run("PUT /es/:target/_doc/:id", func(t *testing.T) {
			t.Run("update document with not exist indexName", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/notExistIndex/_doc/1111", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("update document with exist indexName", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/"+indexName+"/_doc/1111", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("create document with exist indexName not exist id", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/"+indexName+"/_doc/notexist1", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("update document with exist indexName and exist id", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/"+indexName+"/_doc/1111", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("update document with error input", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(`xxx`)
				resp := request("PUT", "/es/"+indexName+"/_doc/1111", body)
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			})
		})

		t.Run("DELETE /es/:target/_doc/:id", func(t *testing.T) {
			t.Run("delete document with not exist indexName", func(t *testing.T) {
				resp := request("DELETE", "/es/notExistIndexDelete/_doc/1111", nil)
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			})
			t.Run("delete document with exist indexName not exist id", func(t *testing.T) {
				resp := request("DELETE", "/es/"+indexName+"/_doc/notexist2", nil)
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			})
			t.Run("delete document with exist indexName and exist id", func(t *testing.T) {
				// wait for WAL write to index
				time.Sleep(time.Second)
				resp := request("DELETE", "/es/"+indexName+"/_doc/1111", nil)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
		})

		t.Run("PUT /es/:target/_create/:id", func(t *testing.T) {
			t.Run("update document with not exist indexName", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/notExistIndexCreate1/_create/1111", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("update document with exist indexName", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/"+indexName+"/_create/1111", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("update document with exist indexName not exist id", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/"+indexName+"/_create/notexistCreate", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("update document with exist indexName and exist id", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/"+indexName+"/_create/1111", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("update document with error input", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(`xxx`)
				resp := request("PUT", "/es/"+indexName+"/_create/1111", body)
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			})
		})

		t.Run("POST /es/:target/_create/:id", func(t *testing.T) {
			t.Run("update document with not exist indexName", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/notExistIndexCreate2/_create/1111", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("update document with exist indexName", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_create/1111", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("update document with exist indexName not exist id", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_create/notexistCreate", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("update document with exist indexName and exist id", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_create/1111", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("update document with error input", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(`xxx`)
				resp := request("POST", "/es/"+indexName+"/_create/1111", body)
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			})
		})

		t.Run("POST /es/:target/_update/:id", func(t *testing.T) {
			t.Run("update document with not exist indexName", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/notExistIndexCreate3/_update/1111", body)
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
			})
			t.Run("update document with exist indexName", func(t *testing.T) {
				// wait for WAL write to index
				time.Sleep(time.Second)
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_update/1111", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("update document with exist indexName not exist id", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_update/notexistCreate", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("update document with exist indexName and exist id", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_update/1111", body)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("update document with error input", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(`xxx`)
				resp := request("POST", "/es/"+indexName+"/_update/1111", body)
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			})
		})

	})
}
