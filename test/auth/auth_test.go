package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prabhatsharma/zinc/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApiBase(t *testing.T) {
	Convey("test auth api", t, func() {
		r := test.Server()
		Convey("check auth with auth", func() {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			req.SetBasicAuth(test.Username, test.Password)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
		})
		Convey("check auth with error password", func() {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			req.SetBasicAuth(test.Username, "xxx")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusUnauthorized)
		})
		Convey("check auth without auth", func() {
			req, _ := http.NewRequest("GET", "/api/index", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusUnauthorized)
		})
	})
}
