package zutils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEnv(t *testing.T) {
	Convey("zutils:env", t, func() {
		Convey("GetEnv", func() {
			a := GetEnv("ZINC_SENTRY", "true")
			So(a, ShouldEqual, "true")
			a = GetEnv("ZINC_SENTRY", "")
			So(a, ShouldEqual, "")
		})
		Convey("GetEnvToUpper", func() {
			a := GetEnvToUpper("ZINC_SENTRY", "true")
			So(a, ShouldEqual, "TRUE")
		})
		Convey("GetEnvToLower", func() {
			a := GetEnvToLower("ZINC_SENTRY", "TRUE")
			So(a, ShouldEqual, "true")
		})
		Convey("GetEnvBool", func() {
			a := GetEnvToBool("ZINC_SENTRY", "true")
			So(a, ShouldEqual, true)
			a = GetEnvToBool("ZINC_SENTRY", "")
			So(a, ShouldEqual, false)
		})
	})
}
