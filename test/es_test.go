package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestApiES(t *testing.T) {
	Convey("test es api", t, func() {
		Convey("POST /es/_bulk", func() {
			Convey("bulk documents", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(bulkData)
				resp := request("POST", "/es/_bulk", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("bulk documents with delete", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(bulkDataWithDelete)
				resp := request("POST", "/es/_bulk", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("bulk with error input", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"index":{}}`)
				resp := request("POST", "/es/_bulk", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})

		Convey("POST /es/:target/_bulk", func() {
			Convey("bulk create documents with not exist indexName", func() {
				body := bytes.NewBuffer(nil)
				data := strings.ReplaceAll(bulkData, `"_index": "games3"`, `"_index": ""`)
				body.WriteString(data)
				resp := request("POST", "/es/notExistIndex/_bulk", body)
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
				resp = request("POST", "/es/"+indexName+"/_bulk", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("bulk with error input", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`{"index":{}}`)
				resp := request("POST", "/es/"+indexName+"/_bulk", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})

		Convey("POST /es/:target/_doc", func() {
			_id := ""
			Convey("create document with not exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/notExistIndex/_doc", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("create document with exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_doc", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("create document with exist indexName not exist id", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_doc", body)
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
				resp := request("POST", "/es/"+indexName+"/_doc", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("create document with error input", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`data`)
				resp := request("POST", "/es/"+indexName+"/_doc", body)
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("PUT /es/:target/_doc/:id", func() {
			Convey("update document with not exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/notExistIndex/_doc/1111", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/"+indexName+"/_doc/1111", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("create document with exist indexName not exist id", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/"+indexName+"/_doc/notexist", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with exist indexName and exist id", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/"+indexName+"/_doc/1111", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with error input", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`xxx`)
				resp := request("PUT", "/es/"+indexName+"/_doc/1111", body)
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("DELETE /es/:target/_doc/:id", func() {
			Convey("delete document with not exist indexName", func() {
				resp := request("DELETE", "/es/notExistIndexDelete/_doc/1111", nil)
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
			Convey("delete document with exist indexName not exist id", func() {
				resp := request("DELETE", "/es/"+indexName+"/_doc/notexist", nil)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("delete document with exist indexName and exist id", func() {
				resp := request("DELETE", "/es/"+indexName+"/_doc/1111", nil)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})

		Convey("PUT /es/:target/_create/:id", func() {
			Convey("update document with not exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/notExistIndexCreate/_create/1111", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/"+indexName+"/_create/1111", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with exist indexName not exist id", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/"+indexName+"/_create/notexistCreate", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with exist indexName and exist id", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("PUT", "/es/"+indexName+"/_create/1111", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with error input", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`xxx`)
				resp := request("PUT", "/es/"+indexName+"/_create/1111", body)
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("POST /es/:target/_create/:id", func() {
			Convey("update document with not exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/notExistIndexCreate/_create/1111", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_create/1111", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with exist indexName not exist id", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_create/notexistCreate", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with exist indexName and exist id", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_create/1111", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with error input", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`xxx`)
				resp := request("POST", "/es/"+indexName+"/_create/1111", body)
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("POST /es/:target/_update/:id", func() {
			Convey("update document with not exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/notExistIndexCreate/_update/1111", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with exist indexName", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_update/1111", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with exist indexName not exist id", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_update/notexistCreate", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with exist indexName and exist id", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(indexData)
				resp := request("POST", "/es/"+indexName+"/_update/1111", body)
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
			Convey("update document with error input", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(`xxx`)
				resp := request("POST", "/es/"+indexName+"/_update/1111", body)
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

	})
}
