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
	"sort"
	"sync"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/metadata"
	"github.com/zinclabs/zinc/pkg/upgrade"
	zincanalysis "github.com/zinclabs/zinc/pkg/uquery/analysis"
	"github.com/zinclabs/zinc/pkg/zutils/hash/rendezvous"
)

var ZINC_INDEX_LIST IndexList

type IndexList struct {
	Indexes map[string]*Index
	lock    sync.RWMutex
}

func SetupIndex() {
	// check version
	version, _ := metadata.KV.Get("version")
	if version == nil {
		// version have version from v0.2.5
		// so if no version, it should be <= v0.2.4
		version = []byte("v0.2.4")
	}

	// start loading index
	ZINC_INDEX_LIST.Indexes = make(map[string]*Index)
	indexes, err := metadata.Cluster.ListIndex(0, 0)
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading index")
	}
	for indexName, metaVersion := range indexes {
		err = LoadIndex(indexName, string(version))
		if err != nil {
			log.Fatal().Err(err).Str("index", indexName).Msg("Error loading index")
		}
		ZINC_CLUSTER.SetIndex(indexName, metaVersion, false)
	}

	// update version
	if string(version) != meta.Version {
		err := metadata.KV.Set("version", []byte(meta.Version))
		if err != nil {
			log.Error().Err(err).Msg("Error set version")
		}
	}
}

func LoadIndex(indexName string, version string) error {
	log.Debug().Str("index", indexName).Str("version", version).Msg("Load index")

	lock, err := metadata.Cluster.NewLocker("index/" + indexName)
	if err != nil {
		return err
	}
	lock.Lock()
	defer lock.Unlock()

	metaIndex, err := metadata.Index.Get(indexName)
	if err != nil {
		return err
	}

	index := new(Index)
	index.ref = metaIndex

	// upgrade from old version
	if version == "" {
		version = meta.Version
	}
	if metaIndex.Meta.GetVersion() != "" {
		version = metaIndex.Meta.GetVersion()
	}
	if version != meta.Version {
		log.Info().Msgf("Upgrade index[%s] from version[%s] to version[%s]", metaIndex.Meta.GetName(), version, meta.Version)
		if err := upgrade.Do(version, metaIndex); err != nil {
			return err
		}
		metaIndex.Meta.SetVersion(meta.Version)
		err := metadata.Index.SetMeta(indexName, metaIndex.Meta.Copy())
		if err != nil {
			return err
		}
	}

	// init shards wrapper
	totalShardNum := 0
	index.shards = make(map[string]*IndexShard, metaIndex.GetShardNum())
	index.localShards = make(map[string]*IndexShard, metaIndex.GetShardNum())
	for _, shard := range metaIndex.Shards.List() {
		id := shard.GetID()
		index.shards[id] = &IndexShard{
			root: index,
			ref:  shard,
			name: index.GetName() + "/" + id,
		}
		index.shards[id].shards = make([]*IndexSecondShard, shard.GetShardNum())
		for i, secondShard := range shard.List() {
			index.shards[id].shards[i] = &IndexSecondShard{
				root: index,
				ref:  secondShard,
			}
			totalShardNum++
		}
	}

	// init shards hashing
	index.shardHashing = rendezvous.New()

	log.Info().Msgf("Loading  index... [%s:%s] shards[%d:%d]", index.GetName(), index.GetStorageType(), index.GetShardNum(), totalShardNum)

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

func ReloadIndex(indexName string) error {
	log.Debug().Str("index", indexName).Msg("Reload index")

	lock, err := metadata.Cluster.NewLocker("index/" + indexName)
	if err != nil {
		return err
	}
	lock.Lock()
	defer lock.Unlock()

	index, ok := GetIndex(indexName)
	if !ok {
		return errors.ErrIndexNotExists
	}

	// update settings
	newSettings, err := metadata.Index.GetSettings(indexName)
	if err != nil {
		return err
	}
	index.SetSettings(newSettings, false)

	// update mappings
	newMappings, err := metadata.Index.GetMappings(indexName)
	if err != nil {
		return err
	}
	index.SetMappings(newMappings, false)

	// update shards
	newShards, err := metadata.Index.GetShards(indexName)
	if err != nil {
		return err
	}
	for _, shard := range newShards.List() {
		err := index.GetRef().Shards.Set(shard)
		if err == nil {
			// new shard
			id := shard.GetID()
			index.shards[id] = &IndexShard{
				root: index,
				ref:  shard,
				name: index.GetName() + "/" + id,
			}
			index.shards[id].shards = make([]*IndexSecondShard, shard.GetShardNum())
			for i, secondShard := range shard.List() {
				index.shards[id].shards[i] = &IndexSecondShard{
					root: index,
					ref:  secondShard,
				}
			}
		}
	}

	// reload index analysis
	ana := index.GetSettings().GetAnalysis()
	if ana != nil {
		index.analyzers, err = zincanalysis.RequestAnalyzer(ana)
		if err != nil {
			return errors.New(errors.ErrorTypeRuntimeException, "parse stored analysis error").Cause(err)
		}
	}

	return nil
}

func (t *IndexList) Set(index *Index) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if _, ok := t.Indexes[index.GetName()]; ok {
		log.Error().Str("index", index.GetName()).Msg("core.IndexList set an exists index")
		return // already exists
	}
	t.Indexes[index.GetName()] = index
}

func (t *IndexList) Get(name string) (*Index, bool) {
	t.lock.RLock()
	idx, ok := t.Indexes[name]
	t.lock.RUnlock()
	return idx, ok
}

func (t *IndexList) GetOrCreate(name, storageType string, shardNum int64) (*Index, bool, error) {
	t.lock.RLock()
	idx, ok := t.Indexes[name]
	t.lock.RUnlock()
	if ok {
		return idx, true, nil
	}

	// local lock
	t.lock.Lock()
	defer t.lock.Unlock()
	// maybe someone else created it while we were waiting for the lock
	idx, ok = t.Indexes[name]
	if ok {
		return idx, true, nil
	}
	// okay, let's create new index
	idx, err := NewIndex(name, storageType, shardNum)
	if err != nil {
		return nil, false, err
	}
	// cache it
	t.Indexes[idx.GetName()] = idx
	return idx, false, nil
}

func (t *IndexList) Delete(name string) {
	t.lock.Lock()
	if idx, ok := t.Indexes[name]; ok {
		if err := idx.Close(); err != nil {
			log.Error().Err(err).Msgf("Error Delete index[%s]", name)
		}
	}
	delete(t.Indexes, name)
	t.lock.Unlock()
}

func (t *IndexList) Len() int {
	t.lock.RLock()
	n := len(t.Indexes)
	t.lock.RUnlock()
	return n
}

func (t *IndexList) List() []*Index {
	t.lock.RLock()
	indexes := make([]*Index, 0, len(t.Indexes))
	for _, index := range t.Indexes {
		indexes = append(indexes, index)
	}
	t.lock.RUnlock()
	return indexes
}

func (t *IndexList) ListStat() []*Index {
	items := t.List()
	return items
}

func (t *IndexList) ListName() []string {
	items := t.List()
	names := make([]string, 0, len(items))
	for _, index := range items {
		names = append(names, index.GetName())
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	return names
}

func (t *IndexList) Close() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	eg := errgroup.Group{}
	eg.SetLimit(config.Global.Shard.GorutineNum)
	for _, index := range t.Indexes {
		index := index
		eg.Go(func() error {
			return index.Close()
		})
	}
	return eg.Wait()
}

// GC auto close unused indexes what inactive for a long time (10m)
func (t *IndexList) GC() error {
	return nil // TODO: implement GC
}
