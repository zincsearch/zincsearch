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

func TestDocumentBulk(t *testing.T) {

	Convey("POST /api/_bulk", t, func() {
		Convey("bulk documents", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(bulkData)
			resp := request("POST", "/api/_bulk", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("bulk documents with delete", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(bulkDataWithDelete)
			resp := request("POST", "/api/_bulk", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("bulk with error input", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"index":{}}`)
			resp := request("POST", "/api/_bulk", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
	})

	Convey("POST /api/:target/_bulk", t, func() {
		Convey("bulk create documents with not exist indexName", func() {
			body := bytes.NewBuffer(nil)
			data := strings.ReplaceAll(bulkData, `"_index": "games3"`, `"_index": ""`)
			body.WriteString(data)
			resp := request("POST", "/api/notExistIndex/_bulk", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("bulk create documents with exist indexName", func() {
			// create index
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"name": "` + indexName + `", "storage_type": "disk"}`)
			resp := request("PUT", "/api/index", body)
			So(resp.Code, ShouldEqual, http.StatusBadRequest)

			respData := make(map[string]string)
			err := json.Unmarshal(resp.Body.Bytes(), &respData)
			So(err, ShouldBeNil)
			So(respData["error"], ShouldEqual, "index ["+indexName+"] already exists")

			// check bulk
			body.Reset()
			data := strings.ReplaceAll(bulkData, `"_index": "games3"`, `"_index": ""`)
			body.WriteString(data)
			resp = request("POST", "/api/"+indexName+"/_bulk", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("bulk with error input", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"index":{}}`)
			resp := request("POST", "/api/"+indexName+"/_bulk", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
	})

}
