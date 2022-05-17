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

	"github.com/goccy/go-json"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDocument(t *testing.T) {

	Convey("PUT /api/:target/_doc", t, func() {
		_id := ""
		Convey("create document with not exist indexName", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/notExistIndex/_doc", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("create document with exist indexName", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/"+indexName+"/_doc", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("create document with exist indexName not exist id", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/"+indexName+"/_doc", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := make(map[string]string)
			err := json.Unmarshal(resp.Body.Bytes(), &data)
			So(err, ShouldBeNil)
			So(data["id"], ShouldNotEqual, "")
			_id = data["id"]
		})
		Convey("update document with exist indexName and exist id", func() {
			body := bytes.NewBuffer(nil)
			data := strings.Replace(indexData, "{", "{\"_id\": \""+_id+"\",", 1)
			body.WriteString(data)
			resp := request("PUT", "/api/"+indexName+"/_doc", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("create document with error input", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`data`)
			resp := request("PUT", "/api/"+indexName+"/_doc", body)
			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("PUT /api/:target/_doc/:id", t, func() {
		Convey("update document with not exist indexName", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/notExistIndex/_doc/1111", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("update document with exist indexName", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/"+indexName+"/_doc/1111", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("create document with exist indexName not exist id", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/"+indexName+"/_doc/notexist", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("update document with exist indexName and exist id", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(indexData)
			resp := request("PUT", "/api/"+indexName+"/_doc/1111", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("update document with error input", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`xxx`)
			resp := request("PUT", "/api/"+indexName+"/_doc/1111", body)
			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("DELETE /api/:target/_doc/:id", t, func() {
		Convey("delete document with not exist indexName", func() {
			resp := request("DELETE", "/api/notExistIndexDelete/_doc/1111", nil)
			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("delete document with exist indexName not exist id", func() {
			resp := request("DELETE", "/api/"+indexName+"/_doc/notexist", nil)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("delete document with exist indexName and exist id", func() {
			resp := request("DELETE", "/api/"+indexName+"/_doc/1111", nil)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
	})

}
