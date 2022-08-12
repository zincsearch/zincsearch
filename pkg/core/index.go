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
	"sync"
	"sync/atomic"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/metadata"
	zincanalysis "github.com/zinclabs/zinc/pkg/uquery/analysis"
	"github.com/zinclabs/zinc/pkg/zutils/hash/rendezvous"
)

type Index struct {
	ref          *meta.Index
	analyzers    map[string]*analysis.Analyzer
	shards       map[string]*IndexShard
	localShards  map[string]*IndexShard
	shardHashing *rendezvous.Rendezvous
	lock         sync.RWMutex
}

func (index *Index) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})
	data["name"] = index.GetName()
	data["storage_type"] = index.GetStorageType()
	data["settings"] = index.GetSettings().Copy()
	data["mappings"] = index.GetMappings().Copy()
	data["stats"] = index.GetStats().Copy()
	data["shard_num"] = index.GetShardNum()
	data["all_shard_num"] = index.GetAllShardNum()
	data["wal_size"] = index.GetWALSize()
	b, err := json.Marshal(data)
	return b, err
}

func (index *Index) ToSimple() *meta.IndexSimple {
	data := new(meta.IndexSimple)
	data.Name = index.GetName()
	data.StorageType = index.GetStorageType()
	data.ShardNum = index.GetShardNum()
	data.AllShardNum = index.GetAllShardNum()
	data.WALSize = index.GetWALSize()
	data.Stats = index.GetStats().Copy()
	data.Settings = index.GetSettings().Copy()
	data.Mappings = make(map[string]interface{})
	for field, prop := range index.GetMappings().ListProperty() {
		data.Mappings[field] = prop.Copy()
	}
	return data
}

func (index *Index) GetName() string {
	return index.ref.GetName()
}

func (index *Index) GetStorageType() string {
	return index.ref.GetStorageType()
}

func (index *Index) GetVersion() string {
	return index.ref.GetStorageType()
}

func (index *Index) GetShardNum() int64 {
	return index.ref.GetShardNum()
}

func (index *Index) GetAllShardNum() int64 {
	var n int64
	for _, shard := range index.shards {
		n += shard.GetShardNum()
	}
	return n
}

func (index *Index) GetRef() *meta.Index {
	return index.ref
}

func (index *Index) GetMeta() *meta.IndexMeta {
	return index.ref.Meta
}

func (index *Index) GetSettings() *meta.IndexSettings {
	return index.ref.Settings
}

func (index *Index) GetMappings() *meta.Mappings {
	return index.ref.Mappings
}

func (index *Index) GetShards() *meta.IndexShards {
	return index.ref.Shards
}

func (index *Index) GetStats() *meta.IndexStat {
	return index.ref.Stats
}

func (index *Index) GetAnalyzers() map[string]*analysis.Analyzer {
	index.lock.RLock()
	a := index.analyzers
	index.lock.RUnlock()
	return a
}

func (index *Index) GetWALSize() uint64 {
	size := uint64(0)
	for _, shard := range index.shards {
		s, err := shard.GetWALSize()
		if err != nil {
			return size
		}
		size += s
	}
	return size
}

func (index *Index) UseTemplate() error {
	template, err := UseTemplate(index.GetName())
	if err != nil {
		return err
	}

	if template == nil {
		return nil
	}

	if template.Template.Settings != nil {
		// update settings
		_ = index.SetSettings(template.Template.Settings, false)
		// update analyzers
		analyzers, _ := zincanalysis.RequestAnalyzer(template.Template.Settings.Analysis)
		index.SetAnalyzers(analyzers)
	}

	if template.Template.Mappings != nil {
		_ = index.SetMappings(template.Template.Mappings, false)
	}

	return nil
}

func (index *Index) SetAnalyzers(analyzers map[string]*analysis.Analyzer) {
	if len(analyzers) == 0 {
		return
	}

	index.lock.Lock()
	index.analyzers = analyzers
	index.lock.Unlock()
}

func (index *Index) SetSettings(settings *meta.IndexSettings, save bool) error {
	log.Debug().Str("index", index.GetName()).Msg("set settings")

	if settings == nil {
		return nil
	}

	// cluster lock
	clusterLock, err := metadata.Cluster.NewLocker("index/settings/" + index.GetName())
	if err != nil {
		return err
	}
	clusterLock.Lock()
	defer clusterLock.Unlock()

	// get latest settings from metadata
	metaSettings, err := metadata.Index.GetSettings(index.GetName())
	if err != nil {
		return err
	}

	// update settings
	if settings.NumberOfShards > 0 {
		metaSettings.SetShards(settings.NumberOfShards)
	}
	if settings.NumberOfReplicas > 0 {
		metaSettings.SetReplicas(settings.NumberOfReplicas)
	}
	if settings.Analysis != nil {
		metaSettings.SetAnalysis(settings.Analysis)
	}
	index.GetRef().Settings.Set(metaSettings)

	if save {
		err := metadata.Index.SetSettings(index.GetName(), index.GetSettings().Copy())
		if err != nil {
			return err
		}
		// notify cluster
		return ZINC_CLUSTER.SetIndex(index.GetName(), time.Now().UnixNano(), true)
	}
	return nil
}

func (index *Index) SetMappings(mappings *meta.Mappings, save bool) error {
	log.Debug().Str("index", index.GetName()).Msg("set mappings")

	if mappings == nil {
		return nil
	}

	// custom analyzer just for text field
	for field, prop := range mappings.ListProperty() {
		if prop.Type != "text" {
			prop.Analyzer = ""
			prop.SearchAnalyzer = ""
			mappings.SetProperty(field, prop)
		}
	}

	// set _id field
	mappings.SetProperty("_id", meta.NewProperty("keyword"))

	// set @timestamp field
	fieldTimestamp, exists := mappings.GetProperty(meta.TimeFieldName)
	if !exists {
		fieldTimestamp = meta.NewProperty("date")
	}
	fieldTimestamp.Index = true
	fieldTimestamp.Sortable = true
	fieldTimestamp.Aggregatable = true
	mappings.SetProperty(meta.TimeFieldName, fieldTimestamp)

	// cluster lock
	clusterLock, err := metadata.Cluster.NewLocker("index/mappings/" + index.GetName())
	if err != nil {
		return err
	}
	clusterLock.Lock()
	defer clusterLock.Unlock()

	// get latest settings from metadata
	metaMappings, err := metadata.Index.GetMappings(index.GetName())
	if err != nil {
		return err
	}

	// merge mappings
	for name, prop := range mappings.ListProperty() {
		metaMappings.SetProperty(name, prop)
	}
	// update mappings
	for name, prop := range metaMappings.ListProperty() {
		index.GetMappings().SetProperty(name, prop)
	}

	if save {
		err := metadata.Index.SetMappings(index.GetName(), index.GetMappings().Copy())
		if err != nil {
			return err
		}
		// notify cluster
		return ZINC_CLUSTER.SetIndex(index.GetName(), time.Now().UnixNano(), true)
	}
	return nil
}

// GetReaders return all shard readers
func (index *Index) GetReaders(timeMin, timeMax int64) ([]*bluge.Reader, error) {
	readers := make([]*bluge.Reader, 0)
	for _, shard := range index.shards {
		rs, err := shard.GetReaders(timeMin, timeMax)
		if err != nil {
			return nil, err
		}
		if len(rs) > 0 {
			readers = append(readers, rs...)
		}
	}
	return readers, nil
}

// UpdateMetadata update index metadata, mainly docNum and storageSize
// need merge from all first layer shards
func (index *Index) UpdateMetadata() error {
	var totalDocNum, totalSize uint64
	for id := range index.shards {
		totalDocNum += atomic.LoadUint64(&index.shards[id].ref.Stats.DocNum)
		totalSize += atomic.LoadUint64(&index.shards[id].ref.Stats.StorageSize)
	}

	if totalDocNum > 0 && totalSize > 0 {
		index.lock.Lock()
		atomic.StoreUint64(&index.ref.Stats.DocNum, totalDocNum)
		atomic.StoreUint64(&index.ref.Stats.StorageSize, totalSize)
		index.lock.Unlock()
	}

	return metadata.Index.SetStats(index.GetName(), index.GetStats().Copy())
}

// UpdateMetadataByShard update first layer shard metadata, mainly docNum, storageSize and timeRange
// need merge from all second layer shards
func (index *Index) UpdateMetadataByShard(id string) {
	var totalDocNum, totalSize uint64
	// update docNum and storageSize
	shard := index.shards[id]
	shardNum := shard.GetShardNum()
	for i := int64(0); i < shardNum; i++ {
		index.UpdateStatsBySecondShard(id, i)
		totalDocNum += atomic.LoadUint64(&shard.ref.Shards[i].Stats.DocNum)
		totalSize += atomic.LoadUint64(&shard.ref.Shards[i].Stats.StorageSize)
		// we just keep latest two shards open
		if i+1 < shardNum {
			secondShard := shard.shards[i]
			secondShard.lock.Lock()
			if secondShard.writer != nil {
				_ = secondShard.writer.Close()
				secondShard.writer = nil
			}
			secondShard.lock.Unlock()
		}
	}
	if totalDocNum > 0 && totalSize > 0 {
		index.lock.Lock()
		atomic.StoreUint64(&shard.ref.Stats.DocNum, totalDocNum)
		atomic.StoreUint64(&shard.ref.Stats.StorageSize, totalSize)
		index.lock.Unlock()
	}

	// update latest shard docTime
	secondShard := shard.shards[shard.GetLatestShardID()]
	index.lock.Lock()
	atomic.StoreInt64(&secondShard.ref.Stats.DocTimeMin, atomic.LoadInt64(&shard.ref.Stats.DocTimeMin))
	atomic.StoreInt64(&secondShard.ref.Stats.DocTimeMax, atomic.LoadInt64(&shard.ref.Stats.DocTimeMax))
	index.lock.Unlock()

	// update local storage
	_ = metadata.Index.SetShard(index.GetName(), shard.ref.Copy())
}

// UpdateStatsBySecondShard update second layer shard stats, mainly docNum and storageSize
func (index *Index) UpdateStatsBySecondShard(id string, secondIndex int64) {
	shard := index.shards[id]
	shard.lock.RLock()
	secondShard := shard.shards[secondIndex]
	shard.lock.RUnlock()

	secondShard.lock.RLock()
	w := secondShard.writer
	secondShard.lock.RUnlock()
	if w == nil {
		return
	}

	var docNum, storageSize uint64
	_, storageSize = w.DirectoryStats()
	if r, err := w.Reader(); err == nil {
		if n, err := r.Count(); err == nil {
			docNum = n
		}
		_ = r.Close()
	}

	index.lock.Lock()
	if docNum > 0 {
		atomic.StoreUint64(&secondShard.ref.Stats.DocNum, docNum)
	}
	if storageSize > 0 {
		atomic.StoreUint64(&secondShard.ref.Stats.StorageSize, storageSize)
	}
	index.lock.Unlock()
}

// Reopen just close the index, it will open automatically by trigger
// Deprecated: it will be removed in the future
func (index *Index) Reopen() error {
	return index.Close()
}

func (index *Index) Close() error {
	eg := errgroup.Group{}
	for _, shard := range index.shards {
		shard := shard
		eg.Go(func() error {
			return shard.Close()
		})
	}
	return eg.Wait()
}
