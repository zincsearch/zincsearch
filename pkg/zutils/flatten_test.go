package zutils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnflatten(t *testing.T) {
	Convey("zutils:unflatten", t, func() {
		data := map[string]interface{}{
			"foo.bar.coo": "abc",
			"foo.bar.oxx": "cbd",
			"foo.bcc.xox": "bdc",
		}
		undata, err := Unflatten(data)
		So(err, ShouldBeNil)
		So(len(undata), ShouldEqual, 1)
		So(len(undata["foo"].(map[string]interface{})), ShouldEqual, 2)
		So(len(undata["foo"].(map[string]interface{})["bar"].(map[string]interface{})), ShouldEqual, 2)
		So(undata["foo"].(map[string]interface{})["bar"].(map[string]interface{})["coo"], ShouldEqual, "abc")
	})
}
