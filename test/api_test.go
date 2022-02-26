package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/prabhatsharma/zinc/pkg/auth"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v1"
	. "github.com/smartystreets/goconvey/convey"
)

type userLoginResponse struct {
	User      auth.ZincUser `json:"user"`
	Validated bool          `json:"validated"`
}

func TestApiStandard(t *testing.T) {
	Convey("test zinc api", t, func() {
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
			Convey("with error username or password", func() {
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

		Convey("DELETE /api/user/:userID", func() {
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

		Convey("GET /api/users", func() {
			resp := request("GET", "/api/users", nil)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			So(err, ShouldBeNil)
			So(data.Hits.Total.Value, ShouldEqual, 1)
			So(data.Hits.Hits[0].ID, ShouldEqual, "admin")
		})

		Convey("PUT /api/index", func() {
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
				resp := request("PUT", "/api/"+indexName+"/_mapping", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
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

		Convey("GET /api/index", func() {
			resp := request("GET", "/api/index", nil)
			So(resp.Code, ShouldEqual, http.StatusOK)

			data := make(map[string]interface{})
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

		Convey("POST /api/_bulk", func() {
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

		Convey("POST /api/:target/_bulk", func() {
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

		Convey("PUT /api/:target/document", func() {
			_id := ""
			Convey("create document with not exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/api/notExistIndex/document", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("create document with exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/api/"+indexName+"/document", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("create document with exist indexName not exist id", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/api/"+indexName+"/document", body)
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
				resp := request("PUT", "/api/"+indexName+"/document", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("create document with error input", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`data`)
				resp := request("PUT", "/api/"+indexName+"/document", body)
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("POST /api/:target/_doc", func() {
			_id := ""
			Convey("create document with not exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/api/notExistIndex/_doc", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("create document with exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/api/"+indexName+"/_doc", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("create document with exist indexName not exist id", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/api/"+indexName+"/_doc", body)
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
				resp := request("POST", "/api/"+indexName+"/_doc", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("create document with error input", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`data`)
				resp := request("POST", "/api/"+indexName+"/_doc", body)
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("PUT /api/:target/_doc/:id", func() {
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

		Convey("DELETE /api/:target/_doc/:id", func() {
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

		Convey("POST /api/:target/_search", func() {
			Convey("init data for search", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/api/"+indexName+"/_doc", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("search document with not exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{}`)
				resp := request("POST", "/api/notExistSearch/_search", body)
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
			Convey("search document with exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"search_type": "alldocuments"}`)
				resp := request("POST", "/api/"+indexName+"/_search", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("search document with not exist term", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"search_type": "match", "query": {"term": "xxxx"}}`)
				resp := request("POST", "/api/"+indexName+"/_search", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(meta.SearchResponse)
				err := json.Unmarshal(resp.Body.Bytes(), data)
				So(err, ShouldBeNil)
				So(data.Hits.Total.Value, ShouldEqual, 0)
			})
			Convey("search document with exist term", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"search_type": "match", "query": {"term": "DEMTSCHENKO"}}`)
				resp := request("POST", "/api/"+indexName+"/_search", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(meta.SearchResponse)
				err := json.Unmarshal(resp.Body.Bytes(), data)
				So(err, ShouldBeNil)
				So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
			})
			Convey("search document type: alldocuments", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"search_type": "alldocuments", "query": {}}`)
				resp := request("POST", "/api/"+indexName+"/_search", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(meta.SearchResponse)
				err := json.Unmarshal(resp.Body.Bytes(), data)
				So(err, ShouldBeNil)
				So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
			})
			Convey("search document type: wildcard", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"search_type": "wildcard", "query": {"term": "dem*"}}`)
				resp := request("POST", "/api/"+indexName+"/_search", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(meta.SearchResponse)
				err := json.Unmarshal(resp.Body.Bytes(), data)
				So(err, ShouldBeNil)
				So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
			})
			Convey("search document type: fuzzy", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"search_type": "fuzzy", "query": {"term": "demtschenk"}}`)
				resp := request("POST", "/api/"+indexName+"/_search", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(meta.SearchResponse)
				err := json.Unmarshal(resp.Body.Bytes(), data)
				So(err, ShouldBeNil)
				So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
			})
			Convey("search document type: term", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{
					"search_type": "term", 
					"query": {
						"term": "Turin", 
						"field":"City"
					}
				}`)
				resp := request("POST", "/api/"+indexName+"/_search", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(meta.SearchResponse)
				err := json.Unmarshal(resp.Body.Bytes(), data)
				So(err, ShouldBeNil)
				So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
			})
			Convey("search document type: daterange", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(fmt.Sprintf(`{
					"search_type": "daterange",
					"query": {
						"start_time": "%s",
						"end_time": "%s"
					}
				}`,
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
			Convey("search document type: matchall", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"search_type": "matchall", "query": {"term": "demtschenk"}}`)
				resp := request("POST", "/api/"+indexName+"/_search", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(meta.SearchResponse)
				err := json.Unmarshal(resp.Body.Bytes(), data)
				So(err, ShouldBeNil)
				So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
			})
			Convey("search document type: match", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"search_type": "match", "query": {"term": "DEMTSCHENKO"}}`)
				resp := request("POST", "/api/"+indexName+"/_search", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(meta.SearchResponse)
				err := json.Unmarshal(resp.Body.Bytes(), data)
				So(err, ShouldBeNil)
				So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
			})
			Convey("search document type: matchphrase", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"search_type": "matchphrase", "query": {"term": "DEMTSCHENKO"}}`)
				resp := request("POST", "/api/"+indexName+"/_search", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(meta.SearchResponse)
				err := json.Unmarshal(resp.Body.Bytes(), data)
				So(err, ShouldBeNil)
				So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
			})
			Convey("search document type: multiphrase", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{
					"search_type": "multiphrase",
					"query": {
						"terms": [
							["demtschenko"],
							["albert"]
						]
					}
				}`)
				resp := request("POST", "/api/"+indexName+"/_search", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(meta.SearchResponse)
				err := json.Unmarshal(resp.Body.Bytes(), data)
				So(err, ShouldBeNil)
				So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
			})
			Convey("search document type: prefix", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"search_type": "prefix", "query": {"term": "dem"}}`)
				resp := request("POST", "/api/"+indexName+"/_search", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(meta.SearchResponse)
				err := json.Unmarshal(resp.Body.Bytes(), data)
				So(err, ShouldBeNil)
				So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
			})
			Convey("search document type: querystring", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"search_type": "querystring", "query": {"term": "DEMTSCHENKO"}}`)
				resp := request("POST", "/api/"+indexName+"/_search", body)
				So(resp.Code, ShouldEqual, http.StatusOK)

				data := new(meta.SearchResponse)
				err := json.Unmarshal(resp.Body.Bytes(), data)
				So(err, ShouldBeNil)
				So(data.Hits.Total.Value, ShouldBeGreaterThanOrEqualTo, 1)
			})
		})

		Convey("POST /api/:target/_search with aggregations", func() {
			Convey("terms aggregation", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{
					"search_type": "matchall", 
					"aggs": {
						"my-agg": {
							"agg_type": "terms",
							"field": "City"
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
					"search_type": "matchall", 
					"aggs": {
						"my-agg-max": {
							"agg_type": "max",
							"field": "Year"
						},
						"my-agg-min": {
							"agg_type": "min",
							"field": "Year"
						},
						"my-agg-avg": {
							"agg_type": "avg",
							"field": "Year"
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
	})
}
