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

	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zinc/pkg/meta"
)

func TestLoadIndexes(t *testing.T) {
	t.Run("load system index", func(t *testing.T) {
		// index cann't be reopen, so need close first
		for _, index := range ZINC_SYSTEM_INDEX_LIST {
			index.Writer.Close()
		}
		var err error
		ZINC_SYSTEM_INDEX_LIST, err = LoadZincSystemIndexes()
		assert.NoError(t, err)
		assert.Equal(t, len(systemIndexList), len(ZINC_SYSTEM_INDEX_LIST))
		assert.Equal(t, "_index_mapping", ZINC_SYSTEM_INDEX_LIST["_index_mapping"].Name)
	})

	t.Run("create some index", func(t *testing.T) {
		index, err := NewIndex("TestLoadIndexes.index_1", "disk", nil)
		assert.NoError(t, err)
		assert.NotNil(t, index)

		err = index.SetSettings(&meta.IndexSettings{
			Analysis: &meta.IndexAnalysis{
				Analyzer: map[string]*meta.Analyzer{
					"default": {
						Type: "standard",
					},
				},
			},
		})
		assert.NoError(t, err)

		err = StoreIndex(index)
		assert.NoError(t, err)
	})

	t.Run("load user index from disk", func(t *testing.T) {
		for _, index := range ZINC_INDEX_LIST {
			index.Writer.Close()
		}
		var err error
		ZINC_INDEX_LIST, err = LoadZincIndexesFromMeta()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(ZINC_INDEX_LIST), 0)
	})

	t.Run("cleanup", func(t *testing.T) {
		err := DeleteIndex("TestLoadIndexes.index_1")
		assert.NoError(t, err)
	})
}
