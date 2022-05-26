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

	"github.com/goccy/go-json"
	. "github.com/smartystreets/goconvey/convey"
)

func TestIndex(t *testing.T) {
	Convey("PUT /api/index", t, func() {
		Convey("create index with payload", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(fmt.Sprintf(`{"name":"%s","storage_type":"disk"}`, "newindex"))
			resp := request("PUT", "/api/index", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
			So(resp.Body.String(), ShouldEqual, `{"index":"newindex","message":"index created","storage_type":"disk"}`)
		})

		Convey("create index with error input", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(fmt.Sprintf(`{"name":"%s","storage_type":"disk"}`, ""))
			resp := request("PUT", "/api/index", body)
			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("GET /api/index", func() {
			resp := request("GET", "/api/index", nil)
			So(resp.Code, ShouldEqual, http.StatusOK)

			// data := make(map[string]interface{})
			data := []interface{}{}
			err := json.Unmarshal(resp.Body.Bytes(), &data)
			So(err, ShouldBeNil)
			So(len(data), ShouldBeGreaterThanOrEqualTo, 1)
		})

		Convey("DELETE /api/index/:indexName", func() {
			Convey("delete index with exist indexName", func() {
				resp := request("DELETE", "/api/index/newindex", nil)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("delete index with not exist indexName", func() {
				resp := request("DELETE", "/api/index/newindex", nil)
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("PUT /api/:target/_mapping", func() {
			Convey("update mappings for index", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{
					"mappings": {
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
					}
				}`)
				resp := request("PUT", "/api/"+indexName+"-mapping/_mapping", body)
				// So(resp.Code, ShouldEqual, http.StatusOK)
				So(resp.Body.String(), ShouldEqual, `{"message":"ok"}`)
			})
		})

		Convey("GET /api/:target/_mapping", func() {
			Convey("get mappings from index", func() {
				resp := request("GET", "/api/"+indexName+"/_mapping", nil)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := make(map[string]interface{})
				err := json.Unmarshal(resp.Body.Bytes(), &data)
				So(err, ShouldBeNil)
				So(data[indexName], ShouldNotBeNil)
				v, ok := data[indexName].(map[string]interface{})
				So(ok, ShouldBeTrue)
				So(v["mappings"], ShouldNotBeNil)
			})
		})
	})
}
