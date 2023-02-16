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
	"time"

	blugeindex "github.com/blugelabs/bluge/index"
	"github.com/stretchr/testify/assert"
	"github.com/zinclabs/zincsearch/pkg/meta"
)

func Test_walMergeDocs_WriteTo(t *testing.T) {
	testData := []map[string]interface{}{
		{
			meta.IDFieldName:     "1",
			meta.ActionFieldName: meta.ActionTypeInsert,
			meta.ShardFieldName:  float64(0),
			meta.TimeFieldName:   float64(time.Now().Unix()),
			"name":               "test",
			"age":                float64(10),
		}, {
			meta.IDFieldName:     "2",
			meta.ActionFieldName: meta.ActionTypeInsert,
			meta.ShardFieldName:  float64(0),
			meta.TimeFieldName:   float64(time.Now().Unix()),
			"name":               "test",
			"age":                float64(10),
		},
	}

	var index *Index
	var shard *IndexShard
	var err error
	docs := make(walMergeDocs)
	t.Run("prepare", func(t *testing.T) {
		for _, d := range testData {
			docs.AddDocument(d)
		}

		index, err = NewIndex("Test_walMergeDocs_WriteTo.index_1", "", 2)
		assert.NoError(t, err)
		assert.NotNil(t, index)

		shard = index.GetShardByDocID("1")
		assert.NotNil(t, shard)

		mappings := meta.NewMappings()
		mappings.SetProperty("name", meta.NewProperty("text"))
		mappings.SetProperty("age", meta.NewProperty("numeric"))
		err = index.SetMappings(mappings)
		assert.NoError(t, err)

		err = StoreIndex(index)
		assert.NoError(t, err)
	})

	t.Run("writeTo", func(t *testing.T) {
		batch := blugeindex.NewBatch()
		err := docs.WriteTo(shard, batch, false)
		assert.NoError(t, err)
	})

	t.Run("rollback", func(t *testing.T) {
		batch := blugeindex.NewBatch()
		err := docs.WriteTo(shard, batch, true)
		assert.NoError(t, err)
	})

	t.Run("Cleanup", func(t *testing.T) {
		err := DeleteIndex("Test_walMergeDocs_WriteTo.index_1")
		assert.NoError(t, err)
	})
}
