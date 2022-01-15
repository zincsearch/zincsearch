package zinc

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetFrontendAssets(t *testing.T) {
	Convey("embed::GetFrontendAssets", t, func() {
		f, err := GetFrontendAssets()
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
		Convey("index.html", func() {
			ff, err := f.Open("index.html")
			So(err, ShouldBeNil)
			fs, err := ff.Stat()
			So(err, ShouldBeNil)
			So(fs.Name(), ShouldEqual, "index.html")
		})
		Convey("manifest.json", func() {
			ff, err := f.Open("manifest.json")
			So(err, ShouldBeNil)
			fs, err := ff.Stat()
			So(err, ShouldBeNil)
			So(fs.Name(), ShouldEqual, "manifest.json")
		})
	})
}
