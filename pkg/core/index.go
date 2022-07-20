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
	"github.com/goccy/go-json"

	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/metadata"
	zincanalysis "github.com/zinclabs/zinc/pkg/uquery/analysis"
)

type Index struct {
	meta.Index
	Analyzers map[string]*analysis.Analyzer `json:"-"`
	lock      sync.RWMutex                  `json:"-"`
	open      uint32                        `json:"-"`
	close     chan struct{}                 `json:"-"`
}

func (index *Index) MarshalJSON() ([]byte, error) {
	index.lock.RLock()
	b, err := json.Marshal(index.Index)
	index.lock.RUnlock()
	return b, err
}

func (index *Index) GetName() string {
	return index.Name
}

func (index *Index) GetMappings() *meta.Mappings {
	index.lock.RLock()
	m := index.Mappings
	index.lock.RUnlock()
	return m
}

func (index *Index) GetSettings() *meta.IndexSettings {
	index.lock.RLock()
	s := index.Settings
	index.lock.RUnlock()
	return s
}
func (index *Index) GetAnalyzers() map[string]*analysis.Analyzer {
	index.lock.RLock()
	a := index.Analyzers
	index.lock.RUnlock()
	return a
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
		// update settings
		_ = index.SetSettings(template.Template.Settings)
		// update analyzers
		analyzers, _ := zincanalysis.RequestAnalyzer(template.Template.Settings.Analysis)
		_ = index.SetAnalyzers(analyzers)
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

	index.lock.Lock()
	index.Settings = settings
	index.lock.Unlock()

	return nil
}

func (index *Index) SetAnalyzers(analyzers map[string]*analysis.Analyzer) error {
	if len(analyzers) == 0 {
		return nil
	}

	index.lock.Lock()
	index.Analyzers = analyzers
	index.lock.Unlock()

	return nil
}

func (index *Index) SetMappings(mappings *meta.Mappings) error {
	if mappings == nil {
		mappings = meta.NewMappings()
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
	index.lock.Lock()
	index.Mappings = mappings
	index.lock.Unlock()

	return nil
}

func (index *Index) SetTimestamp(t int64) {
	index.lock.Lock()
	defer index.lock.Unlock()
	if index.DocTimeMin == 0 {
		index.DocTimeMin = t
		index.DocTimeMax = t
		return
	}
	if t < index.DocTimeMin {
		index.DocTimeMin = t
	} else if t > index.DocTimeMax {
		index.DocTimeMax = t
	}
}

func (index *Index) UpdateMetadata() error {
	var totalDocNum, totalSize uint64
	// update docNum and storageSize
	for i := int64(0); i < atomic.LoadInt64(&index.ShardNum); i++ {
		index.UpdateMetadataByShard(i)
	}
	index.lock.Lock()
	defer index.lock.Unlock()
	for i := int64(0); i < atomic.LoadInt64(&index.ShardNum); i++ {
		totalDocNum += atomic.LoadUint64(&index.Shards[i].DocNum)
		totalSize += atomic.LoadUint64(&index.Shards[i].StorageSize)
	}
	if totalDocNum > 0 && totalSize > 0 {
		atomic.StoreUint64(&index.DocNum, totalDocNum)
		atomic.StoreUint64(&index.StorageSize, totalSize)
	}
	// update docTime
	s := index.Shards[index.GetLatestShardID()]
	atomic.StoreInt64(&s.DocTimeMin, atomic.LoadInt64(&index.DocTimeMin))
	atomic.StoreInt64(&s.DocTimeMax, atomic.LoadInt64(&index.DocTimeMax))

	return metadata.Index.Set(index.Name, index.Index)
}

func (index *Index) UpdateMetadataByShard(n int64) {
	index.lock.RLock()
	s := index.Shards[n]
	index.lock.RUnlock()
	s.Lock.RLock()
	w := s.Writer
	s.Lock.RUnlock()
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
		atomic.StoreUint64(&s.DocNum, docNum)
	}
	if storageSize > 0 {
		atomic.StoreUint64(&s.StorageSize, storageSize)
	}
	index.lock.Unlock()
}

func (index *Index) Reopen() error {
	if err := index.Close(); err != nil {
		return err
	}
	if err := index.OpenWAL(); err != nil {
		return err
	}
	if _, err := index.GetWriter(); err != nil {
		return err
	}
	return nil
}

func (index *Index) Close() error {
	if atomic.LoadUint32(&index.open) == 0 {
		return nil
	}
	index.close <- struct{}{}
	atomic.StoreUint32(&index.open, 0)

	index.lock.Lock()
	defer index.lock.Unlock()
	for _, shard := range index.Shards {
		if shard.Writer == nil {
			continue
		}
		if err := shard.Writer.Close(); err != nil {
			return err
		}
		shard.Writer = nil
	}

	if err := index.WAL.Close(); err != nil {
		return err
	}
	index.WAL = nil

	return nil
}
