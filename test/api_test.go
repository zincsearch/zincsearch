package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

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
				_, userNew, _ := auth.GetUser(username)
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
				exist, _, _ := auth.GetUser(username)
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
				So(resp.Body.String(), ShouldEqual, `{"result":"Index: newindex created","storage_type":"disk"}`)
			})
			Convey("create index with error input", func() {
				body := bytes.NewBuffer(nil)
				body.WriteString(fmt.Sprintf(`{"name":"%s","storage_type":"disk"}`, ""))
				resp := request("PUT", "/api/index", body)
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
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
			})
			Convey("delete index with not exist indexName", func() {
			})
			Convey("delete index with error input", func() {
			})
		})

		Convey("POST /api/_bulk", func() {
			Convey("bulk create documents without indexName", func() {
			})
			Convey("bulk create documents with indexName", func() {
			})
			Convey("bulk with error input", func() {
			})
		})

		Convey("POST /api/:target/_bulk", func() {
			Convey("bulk create documents with not exist indexName", func() {
			})
			Convey("bulk create documents with exist indexName", func() {
			})
			Convey("bulk with error input", func() {
			})
		})

		Convey("PUT /api/:target/document", func() {
			Convey("create document with not exist indexName", func() {
			})
			Convey("create document with exist indexName", func() {
			})
			Convey("create document with exist indexName not exist id", func() {
			})
			Convey("create document with exist indexName and exist id", func() {
			})
			Convey("create document with error input", func() {
			})
		})

		Convey("POST /api/:target/_doc", func() {
			Convey("create document with not exist indexName", func() {
			})
			Convey("create document with exist indexName", func() {
			})
			Convey("create document with exist indexName not exist id", func() {
			})
			Convey("create document with exist indexName and exist id", func() {
			})
			Convey("create document with error input", func() {
			})
		})

		Convey("PUT /api/:target/_doc/:id", func() {
			Convey("update document with not exist indexName", func() {
			})
			Convey("update document with exist indexName", func() {
			})
			Convey("update document with exist indexName not exist id", func() {
			})
			Convey("update document with exist indexName and exist id", func() {
			})
			Convey("update document with error input", func() {
			})
		})

		Convey("DELETE /api/:target/_doc/:id", func() {
			Convey("delete document with not exist indexName", func() {
			})
			Convey("delete document with exist indexName", func() {
			})
			Convey("delete document with exist indexName not exist id", func() {
			})
			Convey("delete document with exist indexName and exist id", func() {
			})
			Convey("delete document with error input", func() {
			})
		})

		Convey("POST /api/:target/_search", func() {
			Convey("search document with not exist indexName", func() {
			})
			Convey("search document with exist indexName", func() {
			})
			Convey("search document with not exist term", func() {
			})
			Convey("search document with exist term", func() {
			})
			Convey("search document type: alldocuments", func() {
			})
			Convey("search document type: wildcard", func() {
			})
			Convey("search document type: fuzzy", func() {
			})
			Convey("search document type: term", func() {
			})
			Convey("search document type: daterange", func() {
			})
			Convey("search document type: matchall", func() {
			})
			Convey("search document type: match", func() {
			})
			Convey("search document type: matchphrase", func() {
			})
			Convey("search document type: multiphrase", func() {
			})
			Convey("search document type: prefix", func() {
			})
			Convey("search document type: querystring", func() {
			})
			Convey("search with error input", func() {
			})
		})

	})
}
