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

package upgrade

import (
	"os"
	"path"
	"sort"

	"github.com/zincsearch/zincsearch/pkg/config"
	"github.com/zincsearch/zincsearch/pkg/ider"
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/zutils"
	"github.com/zincsearch/zincsearch/pkg/zutils/json"
)

// UpgradeFromV026T027 upgrades from version v0.2.6
// this is a break update, meta.Index changed shard field, will lost shards data
// upgrade steps:
// -- update index meta
// -- mv     index index_old
// -- create default layer one shards
// -- mv     index_old/000000 index/shard1/000000
// -- mv     index_old/000001 index/shard1/000001
// -- mv     index_old/000002 index/shard1/000002
// -- mv     index_old/000003 index/shard1/000003
func UpgradeFromV026T027(index *meta.Index) error {
	indexName := index.Name
	rootPath := config.Global.DataPath
	if ok, _ := zutils.IsExist(path.Join(rootPath, indexName)); !ok {
		return nil // if index does not exist, skip
	}

	// update metadata
	if len(index.Shards) == 0 {
		newIndex, err := UpgradeMetadataFromV026T027(nil)
		if err != nil {
			return err
		}
		index.ShardNum = newIndex.ShardNum
		index.Shards = newIndex.Shards
		index.Stats.DocNum = newIndex.Stats.DocNum
		index.Stats.StorageSize = newIndex.Stats.StorageSize
	}

	// check if index has been upgraded
	for id := range index.Shards {
		if ok, _ := zutils.IsExist(path.Join(rootPath, indexName, id, "000000")); ok {
			return nil // if index already upgraded, skip
		}
	}

	// mv index index_old
	err := os.Rename(
		path.Join(rootPath, indexName),
		path.Join(rootPath, indexName+"_old"),
	)
	if err != nil {
		return err
	}
	if err := os.Mkdir(path.Join(rootPath, indexName), 0o755); err != nil {
		return err
	}

	// make new shards
	shardNames := make([]string, 0, index.ShardNum)
	for id := range index.Shards {
		if err := os.Mkdir(path.Join(rootPath, indexName, id), 0o755); err != nil {
			return err
		}
		shardNames = append(shardNames, id)
	}
	sort.Slice(shardNames, func(i, j int) bool {
		return shardNames[i] < shardNames[j]
	})
	firstShardName := shardNames[0]

	// update old shards to new shards
	fs, err := os.ReadDir(path.Join(rootPath, indexName+"_old"))
	if err != nil {
		return err
	}
	for _, f := range fs {
		err := os.Rename(
			path.Join(rootPath, indexName+"_old", f.Name()),
			path.Join(rootPath, indexName, firstShardName, f.Name()),
		)
		if err != nil {
			return err
		}
	}

	// delete empty dir
	return os.Remove(path.Join(rootPath, indexName+"_old"))
}

func UpgradeMetadataFromV026T027(data []byte) (*meta.Index, error) {
	idx026 := new(meta.IndexV026)
	if data == nil {
		data = []byte(`{"shard_num":1,"shards":[{"id":0,"doc_num":0}]}`)
	}
	err := json.Unmarshal(data, idx026)
	if err != nil {
		return nil, err
	}

	shardNames := make([]string, config.Global.Shard.Num)
	shards := make(map[string]*meta.IndexShard, config.Global.Shard.Num)
	for i := int64(0); i < config.Global.Shard.Num; i++ {
		node, err := ider.NewNode(int(i))
		if err != nil {
			return nil, err
		}
		id := node.Generate()
		shard := &meta.IndexShard{ID: id, ShardNum: 1}
		for j := int64(0); j < shard.ShardNum; j++ {
			shard.Shards = append(shard.Shards, &meta.IndexSecondShard{ID: j})
		}
		shards[id] = shard
		shardNames[i] = id
	}

	sort.Slice(shardNames, func(i, j int) bool {
		return shardNames[i] < shardNames[j]
	})
	firstShardName := shardNames[0]

	shards[firstShardName].ShardNum = idx026.ShardNum
	shards[firstShardName].Shards = make([]*meta.IndexSecondShard, idx026.ShardNum)
	for i := int64(0); i < idx026.ShardNum; i++ {
		shards[firstShardName].Shards[i] = &meta.IndexSecondShard{
			ID: i,
			Stats: meta.IndexStat{
				DocNum:      idx026.Shards[i].DocNum,
				StorageSize: idx026.Shards[i].StorageSize,
				DocTimeMin:  idx026.Shards[i].DocTimeMin,
				DocTimeMax:  idx026.Shards[i].DocTimeMax,
			},
		}
		shards[firstShardName].Stats.DocNum += idx026.Shards[i].DocNum
		shards[firstShardName].Stats.StorageSize += idx026.Shards[i].StorageSize
	}

	idx026.Shards = nil
	newJson, err := json.Marshal(idx026)
	if err != nil {
		return nil, err
	}

	idx := new(meta.Index)
	err = json.Unmarshal(newJson, idx)
	if err != nil {
		return nil, err
	}
	idx.Shards = shards
	idx.ShardNum = config.Global.Shard.Num
	idx.Stats.DocNum = shards[firstShardName].Stats.DocNum
	idx.Stats.StorageSize = shards[firstShardName].Stats.StorageSize
	return idx, nil
}
