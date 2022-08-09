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
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/metadata"
	"github.com/zinclabs/zinc/pkg/upgrade"
	zincanalysis "github.com/zinclabs/zinc/pkg/uquery/analysis"
	"github.com/zinclabs/zinc/pkg/zutils/hash/rendezvous"
)

func LoadIndexFromMetadata(indexName string, version string) error {
	lock, err := metadata.Cluster.NewLocker("meta/index/" + indexName)
	if err != nil {
		return err
	}
	lock.Lock()
	defer lock.Unlock()

	input, err := metadata.Index.Get(indexName)
	if err != nil {
		return err
	}

	index := new(Index)
	index.ref = new(meta.Index)
	index.ref.Name = input.Name
	index.ref.StorageType = input.StorageType
	index.ref.Settings = input.Settings
	index.ref.Mappings = input.Mappings
	index.ref.Stats = input.Stats

	// upgrade from old version
	if version == "" {
		version = meta.Version
	}
	if input.Version != "" {
		version = input.Version
	}
	if version != meta.Version {
		log.Info().Msgf("Upgrade index[%s] from version[%s] to version[%s]", input.Name, version, meta.Version)
		if err := upgrade.Do(version, input); err != nil {
			return err
		}
		input.Version = meta.Version
		newData, _ := json.Marshal(input)
		err := metadata.Index.Set(indexName, newData)
		if err != nil {
			return err
		}
		// reload data
		return LoadIndexFromMetadata(indexName, meta.Version)
	}

	// init shards
	index.ref.ShardNum = input.ShardNum
	index.ref.Shards = make(map[string]*meta.IndexShard, index.shardNum)
	for id := range input.Shards {
		index.ref.Shards[id] = &meta.IndexShard{
			ID:       input.Shards[id].ID,
			ShardNum: input.Shards[id].ShardNum,
			Stats:    input.Shards[id].Stats,
		}
		index.ref.Shards[id].Shards = make([]*meta.IndexSecondShard, index.ref.Shards[id].ShardNum)
		for j := range input.Shards[id].Shards {
			index.ref.Shards[id].Shards[j] = &meta.IndexSecondShard{
				ID:    input.Shards[id].Shards[j].ID,
				Stats: input.Shards[id].Shards[j].Stats,
			}
		}
	}

	// init shards wrapper
	totalShardNum := 0
	index.shardNum = index.ref.ShardNum
	index.shards = make(map[string]*IndexShard, index.shardNum)
	index.localShards = make(map[string]*IndexShard, index.shardNum)
	for id := range index.ref.Shards {
		index.shards[id] = &IndexShard{
			root: index,
			ref:  index.ref.Shards[id],
			name: index.ref.Name + "/" + index.ref.Shards[id].ID,
		}
		index.shards[id].shards = make([]*IndexSecondShard, index.ref.Shards[id].ShardNum)
		for j := range index.ref.Shards[id].Shards {
			index.shards[id].shards[j] = &IndexSecondShard{
				root: index,
				ref:  index.ref.Shards[id].Shards[j],
			}
			totalShardNum++
		}
	}

	// init shards hashing
	index.shardHashing = rendezvous.New()
	// for id := range index.shards {
	// 	index.shardHashing.Add(id)
	// }

	log.Info().Msgf("Loading  index... [%s:%s] shards[%d:%d]", index.ref.Name, index.ref.StorageType, index.ref.ShardNum, totalShardNum)

	// load index analysis
	if index.ref.Settings != nil && index.ref.Settings.Analysis != nil {
		index.analyzers, err = zincanalysis.RequestAnalyzer(index.ref.Settings.Analysis)
		if err != nil {
			return errors.New(errors.ErrorTypeRuntimeException, "parse stored analysis error").Cause(err)
		}
	}

	// load in memory
	ZINC_INDEX_LIST.Set(index)

	return nil
}
