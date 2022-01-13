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
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusMovedPermanently)
		})
		Convey("/version", func() {
			req, _ := http.NewRequest("GET", "/version", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)

			data := make(map[string]interface{})
			err := json.Unmarshal(w.Body.Bytes(), &data)
			So(err, ShouldBeNil)
			_, ok := data["Version"]
			So(ok, ShouldBeTrue)
		})
		Convey("/healthz", func() {
			req, _ := http.NewRequest("GET", "/healthz", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)

			data := make(map[string]interface{})
			err := json.Unmarshal(w.Body.Bytes(), &data)
			So(err, ShouldBeNil)
			status, ok := data["status"]
			So(ok, ShouldBeTrue)
			So(status, ShouldEqual, "ok")
		})
		Convey("/ui", func() {
			req, _ := http.NewRequest("GET", "/ui/", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
		})
	})
}
