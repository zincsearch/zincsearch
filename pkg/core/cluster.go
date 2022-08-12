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
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/metadata"
	"github.com/zinclabs/zinc/pkg/zutils"
)

var ZINC_CLUSTER *Cluster

type Cluster struct {
	name         string
	nodes        sync.Map // cluster node
	indexes      sync.Map // cluster indexes
	distribution sync.Map // cluster distribution: index -> shard -> nodeID
	localShards  sync.Map // local held shards
}

func SetupCluster() {
	ZINC_CLUSTER = NewCluster()
	err := ZINC_CLUSTER.Join()
	log.Info().Err(err).Msg("Joining  cluster")
}

func NewCluster() *Cluster {
	cluster := &Cluster{
		name: config.Global.Cluster.Name,
	}
	cluster.HandleClusterEvent()
	cluster.HandleInactiveShards()
	return cluster
}

func (c *Cluster) HandleClusterEvent() {
	go c.handleNodeEvent()
	go c.handleIndexEvent()
	go c.handleDistributionEvent()
}

func (c *Cluster) Name() string {
	return c.name
}

func (c *Cluster) Local() *meta.Node {
	return ZINC_NODE.data
}

func (c *Cluster) GetNodes() []*meta.Node {
	nodes := make([]*meta.Node, 0)
	c.nodes.Range(func(_, value interface{}) bool {
		nodes = append(nodes, value.(*meta.Node))
		return true
	})
	return nodes
}

func (c *Cluster) GetNodesNum() int {
	n := 0
	c.nodes.Range(func(_, _ interface{}) bool {
		n++
		return true
	})
	return n
}

func (c *Cluster) GetNode(nodeID int64) *meta.Node {
	v, ok := c.nodes.Load(nodeID)
	if !ok {
		return nil
	}
	return v.(*meta.Node)
}

func (c *Cluster) GetNodeName(nodeID int64) string {
	v, ok := c.nodes.Load(nodeID)
	if !ok {
		return ""
	}
	return v.(*meta.Node).Name
}

// GetDistribution returns the distribution of the index shards
// returns data is map[indexName][shardName]nodeName
func (c *Cluster) GetDistribution() map[string]map[string]string {
	dis := make(map[string]map[string]string)
	c.distribution.Range(func(indexName, disIndex interface{}) bool {
		dis[indexName.(string)] = make(map[string]string)
		disIndex.(*sync.Map).Range(func(shardName, nodeID interface{}) bool {
			dis[indexName.(string)][shardName.(string)] = c.GetNodeName(nodeID.(int64))
			return true
		})
		return true
	})
	return dis
}

func (c *Cluster) handleNodeEvent() {
	events := metadata.Cluster.Watch(meta.ClusterEventTypeNode)
	for e := range events {
		log.Debug().Str("type", meta.StorageEventTypeString[e.Type]).Str("key", string(e.Key)).Str("value", string(e.Value)).Msg("cluster: node event")
		nodeID, _ := strconv.ParseInt(string(e.Key), 10, 64)
		switch e.Type {
		case meta.StorageEventTypePut, meta.StorageEventTypeCreate, meta.StorageEventTypeUpdate:
			node := new(meta.Node)
			_ = json.Unmarshal(e.Value, node)
			c.nodes.Store(nodeID, node)
			// need release some shards for new node
			c.ReleaseNodeShards()
		case meta.StorageEventTypeDelete:
			// no need to release shards for deleted node
			// because it will be deleted from cluster event
			c.nodes.Delete(nodeID)
		}
	}
}

func (c *Cluster) handleIndexEvent() {
	events := metadata.Cluster.Watch(meta.ClusterEventTypeIndex)
	for e := range events {
		log.Debug().Str("type", meta.StorageEventTypeString[e.Type]).Str("key", string(e.Key)).Str("value", string(e.Value)).Msg("cluster: index event")
		indexName := string(e.Key)
		switch e.Type {
		case meta.StorageEventTypePut, meta.StorageEventTypeCreate, meta.StorageEventTypeUpdate:
			metaVersion, _ := strconv.ParseInt(string(e.Value), 10, 64)
			oldVersion, ok := c.indexes.Load(indexName)
			if ok && oldVersion.(int64) == metaVersion {
				continue // no change
			}
			if !ok {
				// not exists, need load
				err := LoadIndex(indexName, meta.Version)
				if err != nil {
					log.Error().Err(err).Str("index", indexName).Int64("version", metaVersion).Msg("cluster: index event, load index")
					continue
				}
			} else {
				// exists, need reload
				// Reload: just update settings, mappings, shards
				err := ReloadIndex(indexName)
				if err != nil {
					log.Error().Err(err).Str("index", indexName).Int64("version", metaVersion).Msg("cluster: index event, update index")
					continue
				}
			}
			c.indexes.Store(indexName, metaVersion)
		case meta.StorageEventTypeDelete:
			// delete from local cache
			err := c.DeleteIndex(indexName)
			if err != nil {
				log.Error().Err(err).Str("index", indexName).Msg("cluster: index event, delete index")
			}
			// delete from local metadata
			ZINC_INDEX_LIST.Delete(indexName)
		}
	}
}

func (c *Cluster) handleDistributionEvent() {
	events := metadata.Cluster.Watch(meta.ClusterEventTypeDistribution)
	for e := range events {
		log.Debug().Str("type", meta.StorageEventTypeString[e.Type]).Str("key", string(e.Key)).Str("value", string(e.Value)).Msg("cluster: distribution event")
		columns := strings.Split(string(e.Key), "/")
		indexName := columns[0]
		shardName := columns[1]
		disIndex, ok := c.distribution.Load(indexName)
		if !ok {
			disIndex = new(sync.Map)
			c.distribution.Store(indexName, disIndex)
		}
		switch e.Type {
		case meta.StorageEventTypePut, meta.StorageEventTypeCreate, meta.StorageEventTypeUpdate:
			nodeID, _ := strconv.ParseInt(string(e.Value), 10, 64)
			disIndex.(*sync.Map).Store(shardName, nodeID)
		case meta.StorageEventTypeDelete:
			disIndex.(*sync.Map).Store(shardName, int64(0))
		}
	}
}

// Join join the cluster
// 1. lock the cluster meta
// 2. get nodes
// 3. add node
// 4. unlock the cluster meta
func (c *Cluster) Join() error {
	// get lock
	lock, err := metadata.Cluster.NewLocker("nodes")
	if err != nil {
		return err
	}
	lock.Lock()
	defer lock.Unlock()

	// get nodes and calculate the new node id
	nodes, err := metadata.Cluster.ListNode(0, 0)
	if err != nil {
		return err
	}
	// node id: range is 0 to 1023 base on snowflake, just can use 1024 nodes
	// we begin from 1, find a empty slot
	var newNodeID int64 = 1
	for _, node := range nodes {
		c.nodes.Store(node.ID, node)
		if newNodeID == node.ID {
			newNodeID++
		}
	}

	// if not cluster mode, node id always set to 1
	if config.Global.ServerMode != meta.ServerModeCluster {
		newNodeID = 1
	}
	ZINC_NODE.SetID(newNodeID)

	// cache local node
	c.nodes.Store(newNodeID, ZINC_NODE.data)

	// add node
	return metadata.Cluster.Join(newNodeID, c.Local())
}

func (c *Cluster) Leave() error {
	return metadata.Cluster.Leave(c.Local().ID)
}

func (c *Cluster) SetIndex(indexName string, metaVersion int64, update bool) error {
	log.Debug().Str("index", indexName).Int64("version", metaVersion).Msg("cluster: set index")

	c.indexes.Store(indexName, metaVersion)
	if update {
		return metadata.Cluster.SetIndex(indexName, metaVersion)
	}
	return nil
}

func (c *Cluster) DeleteIndex(indexName string) error {
	// delete index from cluster
	c.indexes.Delete(indexName)
	_ = metadata.Cluster.DeleteIndex(indexName)

	// release distribution of this node
	disIndex, ok := c.distribution.Load(indexName)
	if ok {
		disIndex.(*sync.Map).Range(func(shardName, _ interface{}) bool {
			_ = metadata.Cluster.ReleaseDistribute(indexName, shardName.(string))
			return true
		})
	}
	c.distribution.Delete(indexName)

	// delete from local shards
	c.localShards.Range(func(shardName, _ interface{}) bool {
		if strings.HasPrefix(shardName.(string), indexName+"/") {
			c.localShards.Delete(shardName)
		}
		return true
	})

	return nil
}

// LoadDistribution load distribution from cluster
func (c *Cluster) LoadDistribution() {
	c.indexes.Range(func(key, _ interface{}) bool {
		indexName := key.(string)
		data, err := metadata.Cluster.ListDistribution(indexName)
		if err != nil {
			log.Error().Err(err).Str("index", indexName).Msg("cluster: load distribution")
			return false
		}
		disIndex, ok := c.distribution.Load(indexName)
		if !ok {
			disIndex = new(sync.Map)
			c.distribution.Store(indexName, disIndex)
		}
		for shard, nodeID := range data {
			if nodeID == c.Local().ID {
				// skip local node, local distribute need more other operations
				nodeID = int64(0)
			}
			disIndex.(*sync.Map).Store(shard, nodeID)
		}
		return true
	})
}

// DistributeIndexShards distribute shards by the index for the node id
func (c *Cluster) DistributeIndexShards(indexName string) error {
	nodeID := c.Local().ID
	index, ok := GetIndex(indexName)
	if !ok {
		return errors.ErrIndexNotExists
	}
	needShards := int(math.Ceil(float64(index.GetShardNum()) / float64(c.GetNodesNum())))

	shards := make(map[string]bool, needShards)
	disIndex, ok := c.distribution.Load(indexName)
	if !ok {
		disIndex = new(sync.Map)
		c.distribution.Store(indexName, disIndex)
	}
	disIndex.(*sync.Map).Range(func(shardName, value interface{}) bool {
		if nodeID == value.(int64) {
			shards[shardName.(string)] = true
		}
		return true
	})
	if len(shards) >= needShards {
		return nil
	}

	// no local cache, fetch from metadata
	lock, err := metadata.Cluster.NewLocker("distribution/" + indexName)
	if err != nil {
		return err
	}
	lock.Lock()
	defer lock.Unlock()

	// check local again
	disIndex.(*sync.Map).Range(func(shardName, value interface{}) bool {
		if nodeID == value.(int64) {
			shards[shardName.(string)] = true
		}
		return true
	})
	if len(shards) >= needShards {
		return nil
	}

	// get distribution from storage
	indexDistribution, err := metadata.Cluster.ListDistribution(indexName)
	if err != nil {
		return err
	}
	for id := range index.shards {
		if _, ok := indexDistribution[id]; !ok {
			indexDistribution[id] = 0
		}
	}

	for shard, nid := range indexDistribution {
		if nid > 0 && nid != nodeID {
			// if the shard is not distributed to this node, store and skip
			disIndex.(*sync.Map).Store(shard, nid)
			continue
		}
		if len(shards) >= needShards {
			break
		}

		// distribute to this node
		if nid != nodeID {
			err := metadata.Cluster.ShardDistribute(indexName, shard, nodeID)
			if err != nil {
				return err
			}
			// reload shard from storage
			if err := index.ReloadShard(shard); err != nil {
				return err
			}
		}
		shards[shard] = true
		// cache distribution
		disIndex.(*sync.Map).Store(shard, nodeID)
		// add to shards
		index.localShards[shard] = index.shards[shard]
		index.shardHashing.Add(shard)
		// add to local shards
		c.localShards.Store(indexName+"/"+shard, time.Now().Unix())
	}

	if len(shards) > 0 {
		return nil
	}

	// at least one shards is distributed to the node
	// need create a new shard for this node
	shard, err := index.NewShard()
	if err != nil {
		return err
	}
	// distribute to this node
	err = metadata.Cluster.ShardDistribute(indexName, shard, nodeID)
	if err != nil {
		return err
	}
	// cache distribution
	disIndex.(*sync.Map).Store(shard, nodeID)
	index.localShards[shard] = index.shards[shard]
	index.shardHashing.Add(shard)

	// add to local shards
	c.localShards.Store(indexName+"/"+shard, time.Now().Unix())

	return nil
}

// ReleaseIndexShards release some shards let node just hold math.Ceil(shards / nodes)
func (c *Cluster) ReleaseIndexShards(indexName string) error {
	index, ok := GetIndex(indexName)
	if !ok {
		return errors.ErrIndexNotExists
	}
	needShards := int(math.Ceil(float64(index.GetShardNum()) / float64(c.GetNodesNum())))

	disIndex, ok := c.distribution.Load(indexName)
	if !ok {
		return nil
	}
	shards := 0
	nodeID := c.Local().ID
	disIndex.(*sync.Map).Range(func(shardName, value interface{}) bool {
		if nodeID != value.(int64) {
			return true
		}

		shards++
		if shards <= needShards {
			return true
		}

		err := c.ReleaseIndexShard(indexName, shardName.(string))
		if err != nil {
			log.Error().Err(err).Str("index", indexName).Str("shard", shardName.(string)).Msg("cluster: release index shard error")
		}
		return true
	})

	return nil
}

// ReleaseIndexShards release some shards let node just hold math.Ceil(shards / nodes)
func (c *Cluster) ReleaseIndexShard(indexName, shardName string) error {
	disIndex, ok := c.distribution.Load(indexName)
	if !ok {
		return nil
	}
	nodeID := c.Local().ID
	nid, ok := disIndex.(*sync.Map).Load(shardName)
	if !ok {
		return nil
	}
	if nodeID != nid.(int64) {
		return nil
	}

	// close shards
	index, ok := GetIndex(indexName)
	if !ok {
		return errors.ErrIndexNotExists
	}
	index.shards[shardName].Close()
	index.shardHashing.Remove(shardName)
	delete(index.localShards, shardName)
	// release distribution
	err := metadata.Cluster.ReleaseDistribute(indexName, shardName)
	if err != nil {
		return err
	}
	// clear cache
	disIndex.(*sync.Map).Store(shardName, int64(0))
	// clear local shards
	c.localShards.Delete(indexName + "/" + shardName)

	return nil
}

// ReleaseNodeShards release some shards for each index, let node just hold math.Ceil(shards / nodes)
func (c *Cluster) ReleaseNodeShards() {
	c.indexes.Range(func(indexName, version interface{}) bool {
		err := c.ReleaseIndexShards(indexName.(string))
		if err != nil {
			log.Error().Err(err).Str("index", indexName.(string)).Msg("cluster: release node shards")
			return false
		}
		return true
	})
}

// HandleInactiveShards close inactive shards after 10 minutes
func (c *Cluster) HandleInactiveShards() {
	interval, err := zutils.ParseDuration(config.Global.IndexMaxIdleTime)
	if err != nil {
		log.Error().Err(err).Str("IndexMaxIdleTime", config.Global.IndexMaxIdleTime).Msg("cluste: handle inactive shards, ParseDuration error")
		interval = 10 * time.Minute
	}
	intervalSecond := int64(interval.Seconds())

	go func() {
		tick := time.NewTicker(interval)
		for range tick.C {
			now := time.Now().Unix()
			c.localShards.Range(func(key, value interface{}) bool {
				columns := strings.Split(key.(string), "/")
				indexName := columns[0]
				shardName := columns[1]
				if now-value.(int64) > intervalSecond {
					err := c.ReleaseIndexShard(indexName, shardName)
					if err != nil {
						log.Error().Err(err).Str("index", indexName).Str("shard", shardName).Msg("cluster: release index shard error")
					}
				}
				return true
			})
		}
	}()
}
