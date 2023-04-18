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
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/pkg/auth"
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/zutils/json"
)

type userLoginResponse struct {
	User      meta.User `json:"user"`
	Validated bool      `json:"validated"`
}

func TestAuth(t *testing.T) {
	t.Run("test auth api", func(t *testing.T) {
		r := server()
		t.Run("check auth with auth", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			req.SetBasicAuth(username, password)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
		t.Run("check auth with error password", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			req.SetBasicAuth(username, "xxx")
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusUnauthorized, resp.Code)
		})
		t.Run("check auth without auth", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusUnauthorized, resp.Code)
		})

		t.Run("POST /api/login", func(t *testing.T) {
			t.Run("with username and password", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(fmt.Sprintf(`{"_id": "%s", "password": "%s"}`, username, password))
				resp := request("POST", "/api/login", body)
				assert.Equal(t, http.StatusOK, resp.Code)

				data := new(userLoginResponse)
				err := json.Unmarshal(resp.Body.Bytes(), &data)
				assert.NoError(t, err)
				assert.True(t, data.Validated)
			})
			t.Run("with bad username or password", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(fmt.Sprintf(`{"_id": "%s", "password": "xxx"}`, username))
				resp := request("POST", "/api/login", body)
				assert.Equal(t, http.StatusOK, resp.Code)

				data := new(userLoginResponse)
				err := json.Unmarshal(resp.Body.Bytes(), &data)
				assert.NoError(t, err)
				assert.False(t, data.Validated)
			})
		})
	})

	t.Run("test user api", func(t *testing.T) {
		t.Run("PUT /api/user", func(t *testing.T) {
			username := "user1"
			password := "123456"
			t.Run("create user with payload", func(t *testing.T) {
				// create user
				body := bytes.NewBuffer(nil)
				body.WriteString(fmt.Sprintf(`{"_id":"%s","name":"%s","password":"%s","role":"admin"}`, username, username, password))
				resp := request("PUT", "/api/user", body)
				assert.Equal(t, http.StatusOK, resp.Code)

				// login check
				body.Reset()
				body.WriteString(fmt.Sprintf(`{"_id":"%s","password":"%s"}`, username, password))
				resp = request("POST", "/api/login", body)
				assert.Equal(t, http.StatusOK, resp.Code)

				data := new(userLoginResponse)
				err := json.Unmarshal(resp.Body.Bytes(), &data)
				assert.NoError(t, err)
				assert.True(t, data.Validated)
			})
			t.Run("update user", func(t *testing.T) {
				// update user
				body := bytes.NewBuffer(nil)
				body.WriteString(fmt.Sprintf(`{"_id":"%s","name":"%s-updated","password":"%s","role":"admin"}`, username, username, password))
				resp := request("PUT", "/api/user", body)
				assert.Equal(t, http.StatusOK, resp.Code)

				// check updated
				userNew, _, _ := auth.GetUser(username)
				assert.Equal(t, fmt.Sprintf("%s-updated", username), userNew.Name)
			})
			t.Run("create user with error input", func(t *testing.T) {
				body := bytes.NewBuffer(nil)
				body.WriteString(`xxx`)
				resp := request("PUT", "/api/user", body)
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			})
		})

		t.Run("DELETE /api/user/:id", func(t *testing.T) {
			t.Run("delete user with exist userid", func(t *testing.T) {
				username := "user1"
				resp := request("DELETE", "/api/user/"+username, nil)
				assert.Equal(t, http.StatusOK, resp.Code)

				// check user exist
				_, exist, _ := auth.GetUser(username)
				assert.False(t, exist)
			})
			t.Run("delete user with not exist userid", func(t *testing.T) {
				resp := request("DELETE", "/api/user/userNotExist", nil)
				assert.Equal(t, http.StatusOK, resp.Code)
			})
		})

		t.Run("GET /api/user", func(t *testing.T) {
			resp := request("GET", "/api/user", nil)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := make([]meta.User, 0)
			err := json.Unmarshal(resp.Body.Bytes(), &data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(data), 1)
			assert.Equal(t, "admin", data[0].ID)
		})
	})
}
