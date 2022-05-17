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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApiBase(t *testing.T) {
	Convey("test base api", t, func() {
		r := server()
		Convey("/", func() {
			req, _ := http.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, http.StatusFound)
		})
		Convey("/version", func() {
			req, _ := http.NewRequest("GET", "/version", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := make(map[string]interface{})
			err := json.Unmarshal(resp.Body.Bytes(), &data)
			So(err, ShouldBeNil)
			_, ok := data["Version"]
			So(ok, ShouldBeTrue)
		})
		Convey("/healthz", func() {
			req, _ := http.NewRequest("GET", "/healthz", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := make(map[string]interface{})
			err := json.Unmarshal(resp.Body.Bytes(), &data)
			So(err, ShouldBeNil)
			status, ok := data["status"]
			So(ok, ShouldBeTrue)
			So(status, ShouldEqual, "ok")
		})
		Convey("/ui", func() {
			req, _ := http.NewRequest("GET", "/ui/", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
	})
}
