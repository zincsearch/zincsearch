// Copyright 2022 Zinc Labs Inc. and Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
