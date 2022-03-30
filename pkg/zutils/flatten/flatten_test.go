package flatten

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFlattern(t *testing.T) {
	Convey("zutils:flatten", t, func() {
		data := map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": map[string]interface{}{
					"oxx": "cbd",
					"xxo": "dba",
				},
				"arr": []interface{}{"a", "b", "c"},
				"arm": []interface{}{
					map[string]interface{}{
						"a1": "b1",
						"a2": "b2",
					},
					map[string]interface{}{
						"b1": "a1",
						"b2": "a2",
					},
				},
			},
		}
		fdata, err := Flatten(data, "")
		So(err, ShouldBeNil)
		So(len(fdata), ShouldEqual, 7)
		So(fdata["foo.bar.oxx"].(string), ShouldEqual, "cbd")
		So(fdata["foo.arm.1.b1"].(string), ShouldEqual, "a1")
		So(len(fdata["foo.arr"].([]interface{})), ShouldEqual, 3)
	})
}
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
