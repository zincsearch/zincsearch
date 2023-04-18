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
	"fmt"
	"regexp"
	"strings"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"

	"github.com/zincsearch/zincsearch/pkg/bluge/directory"
	"github.com/zincsearch/zincsearch/pkg/config"
	"github.com/zincsearch/zincsearch/pkg/ider"
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/metadata"
	"github.com/zincsearch/zincsearch/pkg/zutils/hash/rendezvous"
)

var indexNameRe = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

func CheckIndexName(name string) error {
	if name == "" {
		return fmt.Errorf("index name cannot be empty")
	}
	if strings.HasPrefix(name, "_") {
		return fmt.Errorf("index name cannot start with _")
	}
	if !indexNameRe.Match([]byte(name)) {
		return fmt.Errorf("index name [%s] is invalid, just accept [a-zA-Z0-9_.-]", name)
	}
	return nil
}

// NewIndex creates an instance of a physical zinc index that can be used to store and retrieve data.
func NewIndex(name, storageType string, shardNum int64) (*Index, error) {
	if err := CheckIndexName(name); err != nil {
		return nil, err
	}

	if storageType == "" {
		storageType = "disk"
	}

	if shardNum <= 0 {
		shardNum = config.Global.Shard.Num
	}

	index := new(Index)
	index.ref = new(meta.Index)
	index.ref.Name = name
	index.ref.StorageType = storageType
	index.ref.Version = meta.Version

	// use template
	if err := index.UseTemplate(); err != nil {
		return nil, err
	}
	if index.ref.Settings != nil {
		if index.ref.Settings.NumberOfShards != 0 {
			shardNum = index.ref.Settings.NumberOfShards
		}
	}

	index.shardNum = shardNum
	index.ref.ShardNum = shardNum
	index.ref.Shards = make(map[string]*meta.IndexShard, index.shardNum)
	for i := int64(0); i < index.shardNum; i++ {
		node, err := ider.NewNode(int(i))
		if err != nil {
			return nil, err
		}
		id := node.Generate()
		shard := &meta.IndexShard{ID: id, ShardNum: 1}
		for j := int64(0); j < shard.ShardNum; j++ {
			shard.Shards = append(shard.Shards, &meta.IndexSecondShard{ID: j})
		}
		index.ref.Shards[id] = shard
	}

	// init shards wrapper
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
		}
	}

	// init shards hashing
	index.shardHashing = rendezvous.New()
	for id := range index.shards {
		index.shardHashing.Add(id)
	}

	return index, nil
}

// LoadIndexWriter load the index writer from the storage
func OpenIndexWriter(name string, storageType string, defaultSearchAnalyzer *analysis.Analyzer, timeRange ...int64) (*bluge.Writer, error) {
	cfg := getOpenConfig(name, storageType, defaultSearchAnalyzer, timeRange...)
	return bluge.OpenWriter(cfg)
}

func getOpenConfig(name string, storageType string, defaultSearchAnalyzer *analysis.Analyzer, timeRange ...int64) bluge.Config {
	dataPath := config.Global.DataPath
	cfg := directory.GetDiskConfig(dataPath, name, timeRange...)
	if defaultSearchAnalyzer != nil {
		cfg.DefaultSearchAnalyzer = defaultSearchAnalyzer
	}
	return cfg
}

// storeIndex stores the index to metadata
func StoreIndex(index *Index) error {
	// check index
	checkIndex(index)
	// store index
	if err := storeIndex(index); err != nil {
		return err
	}
	// cache index
	ZINC_INDEX_LIST.Add(index)
	return nil
}

func checkIndex(index *Index) {
	index.lock.Lock()

	if index.ref.Settings == nil {
		index.ref.Settings = new(meta.IndexSettings)
	}
	if index.ref.Mappings == nil {
		// set default mappings
		index.ref.Mappings = meta.NewMappings()
		index.ref.Mappings.SetProperty(meta.TimeFieldName, meta.NewProperty("date"))
	}
	if index.analyzers == nil {
		index.analyzers = make(map[string]*analysis.Analyzer)
	}

	index.lock.Unlock()
}

func storeIndex(index *Index) error {
	data, err := index.MarshalJSON()
	if err != nil {
		return fmt.Errorf("core.storeIndex: index: %s, error: %s", index.ref.Name, err.Error())
	}
	err = metadata.Index.Set(index.GetName(), data)
	if err != nil {
		return fmt.Errorf("core.storeIndex: index: %s, error: %s", index.ref.Name, err.Error())
	}
	return nil
}

func GetIndex(name string) (*Index, bool) {
	return ZINC_INDEX_LIST.Get(name)
}

func GetOrCreateIndex(name, storageType string, shardNum int64) (*Index, bool, error) {
	return ZINC_INDEX_LIST.GetOrCreate(name, storageType, shardNum)
}
