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

	"github.com/goccy/go-json"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/zinclabs/zinc/pkg/meta"
)

func TestSearch(t *testing.T) {

	Convey("init data for search", t, func() {
		body := bytes.NewBuffer(nil)
		body.WriteString(indexData)
		resp := request("PUT", "/api/"+indexName+"/_doc", body)
		So(resp.Code, ShouldEqual, http.StatusOK)
	})

	Convey("POST /api/:target/_search", t, func() {
		Convey("search document with not exist indexName", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{}`)
			resp := request("POST", "/api/notExistSearch/_search", body)
			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("search document with exist indexName", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"query": {"match_all":{}}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("search document with not exist term", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"query": {"match": {"_all": "xxxx"}}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldEqual, 0)
		})
		Convey("search document with exist term", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"query": {"match": {"_all": "DEMTSCHENKO"}}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: match_all", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"query": {"match_all": {}}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: wildcard", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"query": {"wildcard": {"_all": "dem*"}}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: fuzzy", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"query": {"fuzzy": {"Athlete": "demtschenk"}}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: term", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"query": {"term": {"City": "turin"}}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: daterange", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(
				fmt.Sprintf(`{"query": {"range": {"@timestamp": { "gte": "%s", "lt": "%s"}}}}`,
					time.Now().UTC().Add(time.Hour*-24).Format("2006-01-02T15:04:05Z"),
					time.Now().UTC().Format("2006-01-02T15:04:05Z"),
				))
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: match", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"query": {"match": {"_all": "DEMTSCHENKO"}}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: matchphrase", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"query": {"match_phrase": {"_all": "DEMTSCHENKO"}}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: prefix", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"query": {"prefix": {"_all": "dem"}}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
		Convey("search document type: querystring", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"query": {"query_string": {"query": "DEMTSCHENKO"}}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
		})
	})

	Convey("POST /api/:target/_search with aggregations", t, func() {
		Convey("terms aggregation", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{
				"query": {"match_all":{}}, 
				"aggs": {
					"my-agg-term": {
						"terms": {"field": "City"}
					}
				}
			}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(len(data.Aggregations), ShouldBeGreaterThanOrEqualTo, 1)
		})

		Convey("metric aggregation", func() {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{
				"query": {"match_all":{}}, 
				"aggs": {
					"my-agg-max": {
						"max": {"field": "Year"}
					},
					"my-agg-min": {
						"min": {"field": "Year"}
					},
					"my-agg-avg": {
						"avg": {"field": "Year"}
					}
				}
			}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(len(data.Aggregations), ShouldBeGreaterThanOrEqualTo, 1)
		})
	})

	// Convey("cleanup", t, func() {
	// 	resp := request("DELETE", "/api/index/"+indexName, nil)
	// 	So(resp.Code, ShouldEqual, http.StatusOK)
	// })
}
