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

package index

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/test/utils"
)

func TestList(t *testing.T) {
	t.Run("prepare", func(t *testing.T) {
		index, err := core.NewIndex("TestIndexList.index_1", "disk")
		assert.NoError(t, err)
		assert.NotNil(t, index)

		err = core.StoreIndex(index)
		assert.NoError(t, err)
	})

	t.Run("list", func(t *testing.T) {
		c, w := utils.NewGinContext()
		params := map[string]string{
			"page_num":   "0",
			"page_size":  "20",
			"sort_by":    "name",
			"descending": "false",
			"filter":     "",
		}
		utils.SetGinRequestParams(c, params)
		List(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "")
	})

	t.Run("cleanup", func(t *testing.T) {
		err := core.DeleteIndex("TestIndexList.index_1")
		assert.NoError(t, err)
	})
}

func TestIndexNameList(t *testing.T) {
	t.Run("prepare", func(t *testing.T) {
		index, err := core.NewIndex("TestIndexNameList.index_1", "disk")
		assert.NoError(t, err)
		assert.NotNil(t, index)

		err = core.StoreIndex(index)
		assert.NoError(t, err)
	})

	t.Run("list", func(t *testing.T) {
		c, w := utils.NewGinContext()
		IndexNameList(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "TestIndexNameList.index_1")
	})

	t.Run("cleanup", func(t *testing.T) {
		err := core.DeleteIndex("TestIndexNameList.index_1")
		assert.NoError(t, err)
	})
}
