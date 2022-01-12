package base

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prabhatsharma/zinc/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetFrontendAssets(t *testing.T) {
	Convey("test base api", t, func() {
		r := test.Server()
		Convey("/", func() {
			body := bytes.NewBuffer(nil)
			req, _ := http.NewRequest("GET", "/", body)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
		})
		Convey("/version", func() {

		})
		Convey("/healthz", func() {

		})
		Convey("/ui", func() {

		})
	})
}
