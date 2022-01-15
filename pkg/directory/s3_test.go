package directory

import (
	"testing"

	"github.com/blugelabs/bluge/index"
	. "github.com/smartystreets/goconvey/convey"
)

// check s3 directory implemented index.Directory
var _ index.Directory = new(S3Directory)

func TestUnflatten(t *testing.T) {
	Convey("s3:directory", t, func() {
		Convey("new bluge index directory", func() {

		})
		Convey("create index", func() {

		})
		Convey("load index", func() {

		})
		Convey("create document", func() {

		})
		Convey("update document", func() {

		})
		Convey("list documents", func() {

		})
		Convey("delete document", func() {

		})
		Convey("delete index", func() {

		})
	})
}
