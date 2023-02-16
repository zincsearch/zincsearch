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
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/pkg/zutils/json"
)

func TestIndex(t *testing.T) {
	t.Run("PUT /api/index", func(t *testing.T) {
		t.Run("create index with payload", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(fmt.Sprintf(`{"name":"%s","storage_type":"disk"}`, "newindex"))
			resp := request("PUT", "/api/index", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := make(map[string]interface{})
			err := json.Unmarshal(resp.Body.Bytes(), &data)
			assert.NoError(t, err)
			assert.Equal(t, data["index"], "newindex")
			assert.Equal(t, data["storage_type"], "disk")
			assert.Equal(t, data["message"], "ok")
		})

		t.Run("create index with error input", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(fmt.Sprintf(`{"name":"%s","storage_type":"disk"}`, ""))
			resp := request("PUT", "/api/index", body)
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})

		t.Run("GET /api/index", func(t *testing.T) {
			resp := request("GET", "/api/index", nil)
			assert.Equal(t, http.StatusOK, resp.Code)

			// data := make(map[string]interface{})
			data := struct {
				List []interface{} `json:"list"`
				Page meta.Page     `json:"page"`
			}{}
			err := json.Unmarshal(resp.Body.Bytes(), &data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(data.List), 1)
		})

		t.Run("DELETE /api/index/:indexName", func(t *testing.T) {
			t.Run("delete index with exist indexName", func(t *testing.T) {
				resp := request("DELETE", "/api/index/newindex", nil)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
			t.Run("delete index with not exist indexName", func(t *testing.T) {
				resp := request("DELETE", "/api/index/newindex", nil)
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			})
		})

		t.Run("PUT /api/:target/_mapping", func(t *testing.T) {
			t.Run("update mappings for index", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{
					"properties":{
						"Athlete": {"type": "text"},
						"City": {"type": "keyword"},
						"Country": {"type": "keyword"},
						"Discipline": {"type": "text"},
						"Event": {"type": "keyword"},
						"Gender": {"type": "keyword"},
						"Medal": {"type": "keyword"},
						"Season": {"type": "keyword"},
						"Sport": {"type": "keyword"},
						"Year": {"type": "numeric"},
						"Date": {"type": "time"}
					}
				}`)
				resp := request("PUT", "/api/"+indexName+"-mapping/_mapping", body)
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.Contains(t, resp.Body.String(), `"message":"ok"`)
			})
		})

		t.Run("GET /api/:target/_mapping", func(t *testing.T) {
			t.Run("get mappings from index", func(t *testing.T) {
				resp := request("GET", "/api/"+indexName+"/_mapping", nil)
				assert.Equal(t, http.StatusOK, resp.Code)

				data := make(map[string]interface{})
				err := json.Unmarshal(resp.Body.Bytes(), &data)
				assert.NoError(t, err)
				assert.NotNil(t, data[indexName])
				v, ok := data[indexName].(map[string]interface{})
				assert.True(t, ok)
				assert.NotNil(t, v["mappings"])
			})
		})
	})
}
