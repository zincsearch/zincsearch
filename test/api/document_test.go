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

func TestDocument(t *testing.T) {
	t.Run("PUT /api/:target/_doc", func(t *testing.T) {
		_id := ""
		t.Run("create document with not exist indexName", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/notExistIndex/_doc", body)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
		t.Run("create document with exist indexName", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/"+indexName+"/_doc", body)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
		t.Run("create document with exist indexName not exist id", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/"+indexName+"/_doc", body)
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
			resp := request("PUT", "/api/"+indexName+"/_doc", body)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
		t.Run("create document with error input", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`data`)
			resp := request("PUT", "/api/"+indexName+"/_doc", body)
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})
	})

	t.Run("PUT /api/:target/_doc/:id", func(t *testing.T) {
		t.Run("update document with not exist indexName", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/notExistIndex/_doc/1111", body)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
		t.Run("update document with exist indexName", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/"+indexName+"/_doc/1111", body)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
		t.Run("create document with exist indexName not exist id", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/"+indexName+"/_doc/notexist1", body)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
		t.Run("update document with exist indexName and exist id", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/"+indexName+"/_doc/1111", body)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
		t.Run("update document with error input", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`xxx`)
			resp := request("PUT", "/api/"+indexName+"/_doc/1111", body)
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})
		// wait for WAL write to index
		time.Sleep(time.Second)
	})

	t.Run("DELETE /api/:target/_doc/:id", func(t *testing.T) {
		t.Run("delete document with not exist indexName", func(t *testing.T) {
			resp := request("DELETE", "/api/notExistIndexDelete/_doc/1111", nil)
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})
		t.Run("delete document with exist indexName not exist id", func(t *testing.T) {
			resp := request("DELETE", "/api/"+indexName+"/_doc/notexist2", nil)
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})
		t.Run("delete document with exist indexName and exist id", func(t *testing.T) {
			resp := request("DELETE", "/api/"+indexName+"/_doc/1111", nil)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
	})
}
