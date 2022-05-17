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

	"github.com/goccy/go-json"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/zinclabs/zinc/pkg/auth"
	"github.com/zinclabs/zinc/pkg/meta"
)

type userLoginResponse struct {
	User      auth.ZincUser `json:"user"`
	Validated bool          `json:"validated"`
}

func TestAuth(t *testing.T) {
	Convey("test auth api", t, func() {
		r := server()
		Convey("check auth with auth", func() {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			req.SetBasicAuth(username, password)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("check auth with error password", func() {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			req.SetBasicAuth(username, "xxx")
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, http.StatusUnauthorized)
		})
		Convey("check auth without auth", func() {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("POST /api/login", func() {
			Convey("with username and password", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(fmt.Sprintf(`{"_id": "%s", "password": "%s"}`, username, password))
				resp := request("POST", "/api/login", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(userLoginResponse)
				err := json.Unmarshal(resp.Body.Bytes(), &data)
				So(err, ShouldBeNil)
				So(data.Validated, ShouldBeTrue)
			})
			Convey("with bad username or password", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(fmt.Sprintf(`{"_id": "%s", "password": "xxx"}`, username))
				resp := request("POST", "/api/login", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(userLoginResponse)
				err := json.Unmarshal(resp.Body.Bytes(), &data)
				So(err, ShouldBeNil)
				So(data.Validated, ShouldBeFalse)
			})
		})
	})

	Convey("test user api", t, func() {
		Convey("PUT /api/user", func() {
			username := "user1"
			password := "123456"
			Convey("create user with payload", func() {
				// create user
				body := bytes.NewBuffer(nil)
				body.WriteString(fmt.Sprintf(`{"_id":"%s","name":"%s","password":"%s","role":"admin"}`, username, username, password))
				resp := request("PUT", "/api/user", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				// login check
				body.Reset()
				body.WriteString(fmt.Sprintf(`{"_id":"%s","password":"%s"}`, username, password))
				resp = request("POST", "/api/login", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(userLoginResponse)
				err := json.Unmarshal(resp.Body.Bytes(), &data)
				So(err, ShouldBeNil)
				So(data.Validated, ShouldBeTrue)
			})
			Convey("update user", func() {
				// update user
				body := bytes.NewBuffer(nil)
				body.WriteString(fmt.Sprintf(`{"_id":"%s","name":"%s-updated","password":"%s","role":"admin"}`, username, username, password))
				resp := request("PUT", "/api/user", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				// check updated
				userNew, _, _ := auth.GetUser(username)
				So(userNew.Name, ShouldEqual, fmt.Sprintf("%s-updated", username))
			})
			Convey("create user with error input", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`xxx`)
				resp := request("PUT", "/api/user", body)
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("DELETE /api/user/:id", func() {
			Convey("delete user with exist userid", func() {
				username := "user1"
				resp := request("DELETE", "/api/user/"+username, nil)
				So(resp.Code, ShouldEqual, http.StatusOK)

				// check user exist
				_, exist, _ := auth.GetUser(username)
				So(exist, ShouldBeFalse)
			})
			Convey("delete user with not exist userid", func() {
				resp := request("DELETE", "/api/user/userNotExist", nil)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})

		Convey("GET /api/user", func() {
			resp := request("GET", "/api/user", nil)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldEqual, 1)
			So(data.Hits.Hits[0].ID, ShouldEqual, "admin")
		})
	})
}
