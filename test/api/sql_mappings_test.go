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
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/pkg/config"
	"github.com/zincsearch/zincsearch/pkg/core"
)

func TestSQLMappingsAPI(t *testing.T) {
	enabled := config.Global.EnableSQLMappings
	config.Global.EnableSQLMappings = true
	t.Cleanup(func() { config.Global.EnableSQLMappings = enabled })
	for _, family := range []string{"api", "es"} {
		t.Run(family, func(t *testing.T) {
			name := "testsqlmappingsapi_" + family
			request(http.MethodDelete, "/api/index/"+name, nil) // drop leftovers from a previous run
			t.Cleanup(func() { request(http.MethodDelete, "/api/index/"+name, nil) })
			body := `{"mappings":{"sql":"CREATE TABLE users (id BIGINT, name VARCHAR(100), created_at DATETIME(6), metadata JSON)"}}`
			path := "/es/" + name
			if family == "api" {
				path = "/api/index"
				body = fmt.Sprintf(`{"name":%q,"storage_type":"disk",%s`, name, body[1:])
			}
			resp := request(http.MethodPut, path, bytes.NewBufferString(body))
			if !assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String()) {
				return
			}
			idx, exists := core.GetIndex(name)
			if !assert.True(t, exists) {
				return
			}
			p, ok := idx.GetMappings().GetProperty("name")
			assert.True(t, ok)
			assert.Equal(t, "keyword", p.Type)

			path = "/" + family + "/" + name + "/_mapping"
			resp = request(http.MethodPut, path, bytes.NewBufferString(`{"sql":"CREATE TABLE users (bio TEXT)"}`))
			assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
			p, ok = idx.GetMappings().GetProperty("bio")
			assert.True(t, ok)
			assert.Equal(t, "text", p.Type)
			resp = request(http.MethodPut, path, bytes.NewBufferString(`{"properties":{"tag":{"type":"keyword"}}}`))
			assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
			resp = request(http.MethodPut, path, bytes.NewBufferString(`{"sql":"CREATE TABLE users (extra INT); DROP TABLE users"}`))
			assert.Equal(t, http.StatusBadRequest, resp.Code, resp.Body.String())
			_, ok = idx.GetMappings().GetProperty("extra")
			assert.False(t, ok)

			resp = request(http.MethodPut, "/es/"+name+"/_doc/1", bytes.NewBufferString(`{"id":1,"name":"Ada","created_at":"2026-09-07 12:34:56.123456","metadata":{"region":"west"}}`))
			assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
			assert.Eventually(t, func() bool {
				resp = request(http.MethodGet, "/es/"+name+"/_doc/1", nil)
				return resp.Code == http.StatusOK
			}, 5*time.Second, 20*time.Millisecond)
			assert.Contains(t, resp.Body.String(), `"region":"west"`)
		})
	}

	t.Run("template", func(t *testing.T) {
		resp := request(http.MethodPut, "/es/_index_template/testsqlmappingsapi_template", bytes.NewBufferString(`{"index_patterns":["testsqlmappingsapi_template_*"],"template":{"mappings":{"sql":"CREATE TABLE users (id BIGINT, name VARCHAR(100))"}}}`))
		assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		resp = request(http.MethodGet, "/es/_index_template/testsqlmappingsapi_template", nil)
		assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		assert.Contains(t, resp.Body.String(), `"type":"keyword"`)
	})
}
