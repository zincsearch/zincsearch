/* Copyright 2022 Zinc Labs Inc. and Contributors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package core

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildBlugeDocumentFromJSON(t *testing.T) {
	Convey("test build bluge document from json", t, func() {
		Convey("build bluge document from json", func() {

			idx, _ := NewIndex("index1", "disk", 0, nil)
			// var err error
			// var doc *bluge.Document

			doc1 := make(map[string]interface{})
			doc1["id"] = "1"
			doc1["name"] = "test1"
			doc1["age"] = 10
			doc1["address"] = map[string]interface{}{
				"street": "447 Great Mall Dr",
				"city":   "Milpitas",
				"state":  "CA",
				"zip":    "95035",
			}

			_, err := idx.BuildBlugeDocumentFromJSON("1", doc1)
			So(err, ShouldBeNil)
		})
	})

}
