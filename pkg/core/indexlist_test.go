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

func TestIndexList_List(t *testing.T) {
	indexName := "TestIndexList_List.index_1"
	index, exist, err := GetOrCreateIndex(indexName, "", 2)
	assert.NoError(t, err)
	assert.False(t, exist)
	assert.NotNil(t, index)

	rs, err := index.GetReaders(0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, rs)

	got1 := ZINC_INDEX_LIST.List()
	assert.NotNil(t, got1)

	got2 := ZINC_INDEX_LIST.ListName()
	assert.NotNil(t, got2)

	got3 := ZINC_INDEX_LIST.ListStat()
	assert.NotNil(t, got3)

	err = DeleteIndex(indexName)
	assert.NoError(t, err)

	err = ZINC_INDEX_LIST.GC()
	assert.NoError(t, err)
}

func TestLoadIndexes(t *testing.T) {
	t.Run("create some index", func(t *testing.T) {
		index, _, err := GetOrCreateIndex("TestLoadIndexes.index_1", "disk", 2)
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
		}, true)
		assert.NoError(t, err)
	})

	t.Run("load user index from disk", func(t *testing.T) {
		ZINC_INDEX_LIST.Close()
		err := LoadIndex("TestLoadIndexes.index_1", meta.Version)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, ZINC_INDEX_LIST.Len(), 0)
	})

	t.Run("cleanup", func(t *testing.T) {
		err := DeleteIndex("TestLoadIndexes.index_1")
		assert.NoError(t, err)
	})
}
