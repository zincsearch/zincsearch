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
	"github.com/rs/zerolog/log"

	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/metadata"
	"github.com/zincsearch/zincsearch/pkg/upgrade"
	zincanalysis "github.com/zincsearch/zincsearch/pkg/uquery/analysis"
	"github.com/zincsearch/zincsearch/pkg/zutils/hash/rendezvous"
	"github.com/zincsearch/zincsearch/pkg/zutils/json"
)

func LoadZincIndexesFromMetadata(version string) error {
	indexes, err := metadata.Index.List(0, 0)
	if err != nil {
		return err
	}

	for i := range indexes {
		readIndex := indexes[i]
		index := new(Index)
		index.ref = new(meta.Index)
		index.ref.Name = readIndex.Name
		index.ref.StorageType = readIndex.StorageType
		index.ref.Settings = readIndex.Settings
		index.ref.Mappings = readIndex.Mappings
		index.ref.Stats = readIndex.Stats

		// upgrade from old version
		if readIndex.Version != "" {
			version = readIndex.Version
		}
		if version != meta.Version {
			log.Info().Msgf("Upgrade index[%s] from version[%s] to version[%s]", readIndex.Name, version, meta.Version)
			if err := upgrade.Do(version, readIndex); err != nil {
				return err
			}
			// save updated metadata
			readIndex.Version = meta.Version
			data, err := json.Marshal(readIndex)
			if err != nil {
				return err
			}
			if err := metadata.Index.Set(readIndex.Name, data); err != nil {
				return err
			}
		}

		// init shards
		index.ref.ShardNum = readIndex.ShardNum
		index.ref.Shards = make(map[string]*meta.IndexShard, index.shardNum)
		for id := range readIndex.Shards {
			index.ref.Shards[id] = &meta.IndexShard{
				ID:       readIndex.Shards[id].ID,
				ShardNum: readIndex.Shards[id].ShardNum,
				Stats:    readIndex.Shards[id].Stats,
			}
			index.ref.Shards[id].Shards = make([]*meta.IndexSecondShard, index.ref.Shards[id].ShardNum)
			for j := range readIndex.Shards[id].Shards {
				index.ref.Shards[id].Shards[j] = &meta.IndexSecondShard{
					ID:    readIndex.Shards[id].Shards[j].ID,
					Stats: readIndex.Shards[id].Shards[j].Stats,
				}
			}
		}

		// init shards wrapper
		totalShardNum := 0
		index.shardNum = index.ref.ShardNum
		index.shards = make(map[string]*IndexShard, index.shardNum)
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
		for id := range index.shards {
			index.shardHashing.Add(id)
		}

		log.Info().Msgf("Loading  index... [%s:%s] shards[%d:%d]", index.ref.Name, index.ref.StorageType, index.ref.ShardNum, totalShardNum)

		// load index analysis
		if index.ref.Settings != nil && index.ref.Settings.Analysis != nil {
			index.analyzers, err = zincanalysis.RequestAnalyzer(index.ref.Settings.Analysis)
			if err != nil {
				return errors.New(errors.ErrorTypeRuntimeException, "parse stored analysis error").Cause(err)
			}
		}

		// load in memory
		ZINC_INDEX_LIST.Add(index)
	}

	return nil
}
