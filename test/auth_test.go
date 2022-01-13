package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestApiAuth(t *testing.T) {
	Convey("test auth api", t, func() {
		r := server()
		Convey("check auth with auth", func() {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			req.SetBasicAuth(username, password)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
		Convey("check auth with error password", func() {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			req.SetBasicAuth(username, "xxx")
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, http.StatusUnauthorized)
		})
		Convey("check auth without auth", func() {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, http.StatusUnauthorized)
		})
	})
}
