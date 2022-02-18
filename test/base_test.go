package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
