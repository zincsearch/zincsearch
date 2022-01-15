package core

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadIndexes(t *testing.T) {
	Convey("test load index", t, func() {
		Convey("load system index", func() {
			indexes, err := LoadZincSystemIndexes()
			So(err, ShouldBeNil)
			So(len(indexes), ShouldEqual, 2)
			So(indexes["_index_mapping"].Name, ShouldEqual, "_index_mapping")
		})
		Convey("load user inex from disk", func() {
			indexes, err := LoadZincIndexesFromDisk()
			So(err, ShouldBeNil)
			So(len(indexes), ShouldBeGreaterThanOrEqualTo, 0)
		})
		Convey("load user inex from s3", func() {
			// TODO: support
		})
	})
}
