package core

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadIndexes(t *testing.T) {
	Convey("test load index", t, func() {
		Convey("load system index", func() {
			// index cann't be reopen, so need close first
			for _, index := range ZINC_SYSTEM_INDEX_LIST {
				index.Writer.Close()
			}
			var err error
			ZINC_SYSTEM_INDEX_LIST, err = LoadZincSystemIndexes()
			So(err, ShouldBeNil)
			So(len(ZINC_SYSTEM_INDEX_LIST), ShouldEqual, len(systemIndexList))
			So(ZINC_SYSTEM_INDEX_LIST["_index_mapping"].Name, ShouldEqual, "_index_mapping")
		})
		Convey("load user inex from disk", func() {
			// index cann't be reopen, so need close first
			for _, index := range ZINC_INDEX_LIST {
				index.Writer.Close()
			}
			var err error
			ZINC_INDEX_LIST, err = LoadZincIndexesFromDisk()
			So(err, ShouldBeNil)
			So(len(ZINC_INDEX_LIST), ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}
