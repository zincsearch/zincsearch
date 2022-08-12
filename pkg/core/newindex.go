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
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"

	"github.com/zinclabs/zinc/pkg/bluge/directory"
	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/ider"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/metadata"
	"github.com/zinclabs/zinc/pkg/zutils/hash/rendezvous"
)

func GetIndex(name string) (*Index, bool) {
	return ZINC_INDEX_LIST.Get(name)
}

func GetOrCreateIndex(name, storageType string, shardNum int64) (*Index, bool, error) {
	return ZINC_INDEX_LIST.GetOrCreate(name, storageType, shardNum)
}

// NewIndex creates an instance of a physical zinc index that can be used to store and retrieve data.
func NewIndex(name, storageType string, shardNum int64) (*Index, error) {
	// cluster lock
	clusterLock, err := metadata.Cluster.NewLocker("index/" + name)
	if err != nil {
		return nil, err
	}
	clusterLock.Lock()
	defer clusterLock.Unlock()

	// check index name
	if err := checkIndexName(name); err != nil {
		return nil, err
	}

	// check repeat create
	if idx, _ := metadata.Index.GetMeta(name); idx != nil {
		return nil, errors.ErrIndexIsExists
	}

	// create index meta
	if storageType == "" {
		storageType = "disk"
	}
	metaIndex := meta.NewIndex(name, storageType, meta.Version)
	index := new(Index)
	index.ref = metaIndex
	index.analyzers = make(map[string]*analysis.Analyzer)

	// use template
	if err := index.UseTemplate(); err != nil {
		return nil, err
	}

	// check shard num
	if shardNum <= 0 {
		shardNum = config.Global.Shard.Num
	}
	if metaIndex.Settings.GetShards() != 0 {
		shardNum = metaIndex.Settings.GetShards()
	}

	// init shards
	for i := int64(0); i < shardNum; i++ {
		node, err := ider.NewNode(int(i))
		if err != nil {
			return nil, err
		}
		_, err = metaIndex.Shards.Create(node.Generate())
		if err != nil {
			return nil, err
		}
	}

	// init shards wrapper
	index.shards = make(map[string]*IndexShard, shardNum)
	index.localShards = make(map[string]*IndexShard, shardNum)
	for _, shard := range metaIndex.Shards.List() {
		id := shard.GetID()
		index.shards[id] = &IndexShard{
			root: index,
			ref:  shard,
			name: metaIndex.GetName() + "/" + id,
		}
		index.shards[id].shards = make([]*IndexSecondShard, shard.GetShardNum())
		for i, secondShard := range shard.List() {
			index.shards[id].shards[i] = &IndexSecondShard{
				root: index,
				ref:  secondShard,
			}
		}
	}

	// init shards hashing
	index.shardHashing = rendezvous.New()

	// store in local
	err = metadata.Index.SetMeta(name, metaIndex.Meta.Copy())
	if err != nil {
		return nil, err
	}
	err = metadata.Index.SetStats(name, metaIndex.Stats.Copy())
	if err != nil {
		return nil, err
	}
	err = metadata.Index.SetSettings(name, metaIndex.Settings.Copy())
	if err != nil {
		return nil, err
	}
	err = metadata.Index.SetMappings(name, metaIndex.Mappings.Copy())
	if err != nil {
		return nil, err
	}
	err = metadata.Index.SetShards(name, metaIndex.Shards.Copy())
	if err != nil {
		return nil, err
	}

	// notify cluster
	err = ZINC_CLUSTER.SetIndex(name, time.Now().UnixNano(), true)
	if err != nil {
		return nil, err
	}

	return index, nil
}

var indexNameRe = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

func checkIndexName(name string) error {
	if name == "" {
		return errors.ErrIndexIsEmpty
	}
	if strings.HasPrefix(name, "_") {
		return fmt.Errorf("index name cannot start with _")
	}
	if !indexNameRe.Match([]byte(name)) {
		return fmt.Errorf("index name [%s] is invalid, just accept [a-zA-Z0-9_.-]", name)
	}
	return nil
}

// openIndexReader load the index reader from the storage
func openIndexReader(name string, storageType string, ans *analysis.Analyzer, timeRange ...int64) (*bluge.Reader, error) {
	cfg := getOpenConfig(name, storageType, ans, timeRange...)
	return bluge.OpenReader(cfg)
}

// openIndexWriter load the index writer from the storage
func openIndexWriter(name string, storageType string, ans *analysis.Analyzer, timeRange ...int64) (*bluge.Writer, error) {
	fmt.Println("openIndexWriter", name, storageType)

	cfg := getOpenConfig(name, storageType, ans, timeRange...)
	return bluge.OpenWriter(cfg)
}

func getOpenConfig(name string, storageType string, ans *analysis.Analyzer, timeRange ...int64) bluge.Config {
	var dataPath string
	var cfg bluge.Config
	switch storageType {
	case "s3":
		dataPath = config.Global.S3.Bucket
		cfg = directory.GetS3Config(dataPath, name, timeRange...)
	case "minio":
		dataPath = config.Global.MinIO.Bucket
		cfg = directory.GetMinIOConfig(dataPath, name, timeRange...)
	default:
		dataPath = config.Global.DataPath
		cfg = directory.GetDiskConfig(dataPath, name, timeRange...)
	}
	if ans != nil {
		cfg.DefaultSearchAnalyzer = ans
	}
	return cfg
}
