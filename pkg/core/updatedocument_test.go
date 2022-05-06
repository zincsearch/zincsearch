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
	"math/rand"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpdateDocument(t *testing.T) {
	Convey("UpdateDocument", t, func() {
		rand.Seed(time.Now().UnixNano())
		id := rand.Intn(1000)
		indexName := "TestUpdateDocument.index_" + strconv.Itoa(id)

		index, _ := NewIndex(indexName, "disk", UseNewIndexMeta, nil)
		Convey("insert with mintedID", func() {
			err := index.UpdateDocument("doc1", map[string]interface{}{
				"name": "doc1",
			}, true)
			So(err, ShouldBeNil)
			So(index.DocsCount, ShouldEqual, 1)
		})

		Convey("insert with provided ID", func() {
			err := index.UpdateDocument("doc2", map[string]interface{}{
				"Id":   "17",
				"name": "doc1",
			}, false)
			So(err, ShouldBeNil)
			// So(index.DocsCount, ShouldEqual, 1) // TODO: fix this. search for the document with ID 17 should return 1 document
		})

	})

}
