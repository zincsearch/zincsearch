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

	"github.com/blugelabs/bluge/analysis"

	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/metadata"
)

type Index struct {
	meta.Index
	Analyzers map[string]*analysis.Analyzer `json:"-"`
	lock      sync.RWMutex                  `json:"-"`

	wal WriteAheadLog
}

func (index *Index) UseTemplate() error {
	template, err := UseTemplate(index.Name)
	if err != nil {
		return err
	}

	if template == nil {
		return nil
	}

	if template.Template.Settings != nil {
		_ = index.SetSettings(template.Template.Settings)
	}

	if template.Template.Mappings != nil {
		_ = index.SetMappings(template.Template.Mappings)
	}

	return nil
}

func (index *Index) SetSettings(settings *meta.IndexSettings) error {
	if settings == nil {
		return nil
	}

	index.Settings = settings

	return nil
}

func (index *Index) SetAnalyzers(analyzers map[string]*analysis.Analyzer) error {
	if len(analyzers) == 0 {
		return nil
	}

	index.Analyzers = analyzers

	return nil
}

func (index *Index) SetMappings(mappings *meta.Mappings) error {
	if mappings == nil || mappings.Len() == 0 {
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

	mappings.SetProperty("_id", meta.NewProperty("keyword"))

	// @timestamp need date_range/date_histogram aggregation, and mappings used for type check in aggregation
	fieldTimestamp, exists := mappings.GetProperty(meta.TimeFieldName)
	if !exists {
		fieldTimestamp = meta.NewProperty("date")
	}
	fieldTimestamp.Index = true
	fieldTimestamp.Sortable = true
	fieldTimestamp.Aggregatable = true
	mappings.SetProperty(meta.TimeFieldName, fieldTimestamp)

	// update in the cache
	index.Mappings = mappings

	return nil
}

func (index *Index) SetTimestamp(t int64) {
	if index.DocTimeMin == 0 {
		atomic.StoreInt64(&index.DocTimeMin, t)
	}
	if index.DocTimeMax == 0 {
		atomic.StoreInt64(&index.DocTimeMax, t)
	}
	if t < index.DocTimeMin {
		atomic.StoreInt64(&index.DocTimeMin, t)
	}
	if t > index.DocTimeMax {
		atomic.StoreInt64(&index.DocTimeMax, t)
	}
}

func (index *Index) UpdateMetadata() error {
	var totalDocNum, totalSize uint64
	// update docNum and storageSize
	for i := 0; i < index.ShardNum; i++ {
		index.UpdateMetadataByShard(i)
	}
	index.lock.RLock()
	for i := 0; i < index.ShardNum; i++ {
		totalDocNum += index.Shards[i].DocNum
		totalSize += index.Shards[i].StorageSize
	}
	if totalDocNum > 0 && totalSize > 0 {
		index.DocNum = totalDocNum
		index.StorageSize = totalSize
	}
	// update docTime
	index.Shards[index.ShardNum-1].DocTimeMin = index.DocTimeMin
	index.Shards[index.ShardNum-1].DocTimeMax = index.DocTimeMax
	index.lock.RUnlock()

	return metadata.Index.Set(index.Name, index.Index)
}

func (index *Index) UpdateMetadataByShard(n int) {
	index.lock.RLock()
	shard := index.Shards[n]
	index.lock.RUnlock()
	if shard.Writer == nil {
		return
	}
	var docNum, storageSize uint64
	_, storageSize = shard.Writer.DirectoryStats()
	if r, err := shard.Writer.Reader(); err == nil {
		if n, err := r.Count(); err == nil {
			docNum = n
		}
		_ = r.Close()
	}
	if docNum > 0 {
		shard.DocNum = docNum
	}
	if storageSize > 0 {
		shard.StorageSize = storageSize
	}
}

func (index *Index) Reopen() error {
	if err := index.Close(); err != nil {
		return err
	}
	if _, err := index.GetWriter(); err != nil {
		return err
	}
	return nil
}

func (index *Index) Close() error {
	var err error
	// update metadata before close
	if err = index.UpdateMetadata(); err != nil {
		return err
	}

	// TODO: flush WAL or not?

	index.lock.Lock()
	for _, shard := range index.Shards {
		if shard.Writer == nil {
			continue
		}
		if e := shard.Writer.Close(); e != nil {
			err = e
		}
		shard.Writer = nil
	}
	index.lock.Unlock()
	return err
}
