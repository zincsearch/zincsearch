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
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/zincsearch/zincsearch/pkg/zutils/json"
)

func TestApiBase(t *testing.T) {
	t.Run("test base api", func(t *testing.T) {
		r := server()
		t.Run("/", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusFound, resp.Code)
		})
		t.Run("/version", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/version", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := make(map[string]interface{})
			err := json.Unmarshal(resp.Body.Bytes(), &data)
			assert.NoError(t, err)
			_, ok := data["version"]
			assert.True(t, ok)
		})
		t.Run("/healthz", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/healthz", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := make(map[string]interface{})
			err := json.Unmarshal(resp.Body.Bytes(), &data)
			assert.NoError(t, err)
			status, ok := data["status"]
			assert.True(t, ok)
			assert.Equal(t, "ok", status)
		})
		for _, path := range []string{"/ui/search?q=test", "/ui/index", "/ui/template", "/ui/user", "/ui/role", "/ui/about", "/ui/login"} {
			t.Run(path, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, path, nil)
				resp := httptest.NewRecorder()
				r.ServeHTTP(resp, req)
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.Contains(t, resp.Body.String(), `id="root"`)
				assert.Contains(t, resp.Body.String(), `/ui/assets/`)
				assert.Empty(t, resp.Header().Get("Location"))
			})
		}
		t.Run("missing asset", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/ui/assets/missing.js", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusNotFound, resp.Code)
		})
		t.Run("/ui", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/ui/", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
		t.Run("no route logging", func(t *testing.T) {
			var buf bytes.Buffer
			orig := log.Logger
			log.Logger = zerolog.New(&buf)
			defer func() { log.Logger = orig }()

			r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/ui/search", nil))
			assert.NotContains(t, buf.String(), `"code":404`)

			buf.Reset()
			r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/nope", nil))
			assert.Contains(t, buf.String(), `"code":404`)
			assert.Contains(t, buf.String(), "/nope")
		})
	})
}
