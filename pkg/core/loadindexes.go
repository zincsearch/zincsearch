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

	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/metadata"
	"github.com/zinclabs/zinc/pkg/upgrade"
	zincanalysis "github.com/zinclabs/zinc/pkg/uquery/analysis"
)

func LoadZincIndexesFromMetadata() error {
	indexes, err := metadata.Index.List(0, 0)
	if err != nil {
		return err
	}

	for i := range indexes {
		// cache mappings
		index := new(Index)
		index.Name = indexes[i].Name
		index.StorageType = indexes[i].StorageType
		index.StorageSize = indexes[i].StorageSize
		index.DocTimeMin = indexes[i].DocTimeMin
		index.DocTimeMax = indexes[i].DocTimeMax
		index.DocNum = indexes[i].DocNum
		index.ShardNum = indexes[i].ShardNum
		index.Shards = append(index.Shards, indexes[i].Shards...)
		index.Settings = indexes[i].Settings
		index.Mappings = indexes[i].Mappings
		index.CreateAt = indexes[i].CreateAt
		index.UpdateAt = indexes[i].UpdateAt
		index.close = make(chan struct{})

		log.Info().Msgf("Loading  index... [%s:%s] shards[%d]", index.Name, index.StorageType, index.ShardNum)

		// upgrade from version <= 0.2.4
		if index.ShardNum == 0 {
			index.ShardNum = 1
			index.Shards = append(index.Shards, &meta.IndexShard{})
			//upgrade data
			if index.StorageType != "disk" {
				log.Panic().Msgf("Only disk storage type support upgrade from version <= 0.2.4, Please manual upgrade\n# mv %s %s_bak\n# mkdir %s\n# mv %s_bak %s/000000\n# restart zinc", index.Name, index.Name, index.Name, index.Name, index.Name)
			} else {
				if err := upgrade.UpgradeFromV024Index(index.Name); err != nil {
					log.Panic().Err(err).Msgf("Automatic upgrade from version <= 0.2.4 failed, Please manual upgrade\n# mv %s %s_bak\n# mkdir %s\n# mv %s_bak %s/000000\n# restart zinc", index.Name, index.Name, index.Name, index.Name, index.Name)
				}
			}
		}

		// load index analysis
		if index.Settings != nil && index.Settings.Analysis != nil {
			index.Analyzers, err = zincanalysis.RequestAnalyzer(index.Settings.Analysis)
			if err != nil {
				return errors.New(errors.ErrorTypeRuntimeException, "parse stored analysis error").Cause(err)
			}
		}

		// load in memory
		ZINC_INDEX_LIST.Add(index)
	}

	return nil
}
