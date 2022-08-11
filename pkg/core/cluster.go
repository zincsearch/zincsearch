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

	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/metadata"
)

var ZINC_CLUSTER *Cluster

type Cluster struct {
	nodes        sync.Map
	indexes      sync.Map
	distribution sync.Map
}

func SetupCluster() {
	ZINC_CLUSTER = NewCluster()
	err := ZINC_CLUSTER.Join()
	log.Info().Err(err).Msg("Joining  cluster")
}

func NewCluster() *Cluster {
	cluster := &Cluster{}
	cluster.HandleClusterEvent()
	return cluster
}

func (c *Cluster) HandleClusterEvent() {
	go c.handleNodeEvent()
	go c.handleIndexEvent()
	go c.handleDistributionEvent()
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
	c.nodes.Range(func(_, value interface{}) bool {
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
	c.distribution.Range(func(indexName, value interface{}) bool {
		dis[indexName.(string)] = make(map[string]string)
		value.(*sync.Map).Range(func(shardName, value interface{}) bool {
			dis[indexName.(string)][shardName.(string)] = c.GetNodeName(value.(int64))
			return true
		})
		return true
	})
	return dis
}

func (c *Cluster) handleNodeEvent() {
	events := metadata.Cluster.Watch(meta.ClusterEventTypeNode)
	for e := range events {
		nodeID, _ := strconv.ParseInt(string(e.Key), 10, 64)
		switch e.Type {
		case meta.StorageEventTypePut:
			node := new(meta.Node)
			_ = json.Unmarshal(e.Value, node)
			c.nodes.Store(nodeID, node)
			// need release some shards for new node
			c.ReleaseNodeShards(c.Local().ID)
		case meta.StorageEventTypeDelete:
			c.nodes.Delete(nodeID)
		}
		log.Debug().Str("type", meta.StorageEventTypeString[e.Type]).Str("key", string(e.Key)).Str("value", string(e.Value)).Msg("cluster node event")
	}
}

func (c *Cluster) handleIndexEvent() {
	events := metadata.Cluster.Watch(meta.ClusterEventTypeIndex)
	for e := range events {
		indexName := string(e.Key)
		switch e.Type {
		case meta.StorageEventTypePut:
			metaVersion, _ := strconv.ParseInt(string(e.Value), 10, 64)
			oldVersion, ok := c.indexes.Load(indexName)
			if ok && oldVersion.(int64) == metaVersion {
				continue // no change
			}
			if !ok {
				// not exists, need load
				err := LoadIndex(indexName, meta.Version)
				if err != nil {
					log.Error().Err(err).Str("index", indexName).Int64("version", metaVersion).Msg("cluster index event: load index")
					continue
				}
			} else {
				// exists, need reload
				// Reload: just update settings, mappings, shards
				err := ReloadIndex(indexName)
				if err != nil {
					log.Error().Err(err).Str("index", indexName).Int64("version", metaVersion).Msg("cluster index event: update index")
					continue
				}
			}
			c.indexes.Store(indexName, metaVersion)
		case meta.StorageEventTypeDelete:
			// delete from local cache
			err := c.DeleteIndex(indexName)
			if err != nil {
				log.Error().Err(err).Str("index", indexName).Msg("cluster index event, delete index")
			}
			// delete from local metadata
			ZINC_INDEX_LIST.Delete(indexName)
		}
		log.Debug().Str("type", meta.StorageEventTypeString[e.Type]).Str("key", string(e.Key)).Str("value", string(e.Value)).Msg("cluster index event")
	}
}

func (c *Cluster) handleDistributionEvent() {
	events := metadata.Cluster.Watch(meta.ClusterEventTypeDistribution)
	for e := range events {
		columns := strings.Split(string(e.Key), "/")
		indexName := columns[0]
		shardName := columns[1]
		switch e.Type {
		case meta.StorageEventTypePut:
			nodeID, _ := strconv.ParseInt(string(e.Value), 10, 64)
			disIndex, ok := c.distribution.Load(indexName)
			if !ok {
				disIndex = new(sync.Map)
				c.distribution.Store(indexName, disIndex)
			}
			disIndex.(*sync.Map).Store(shardName, nodeID)
		case meta.StorageEventTypeDelete:
			disIndex, ok := c.distribution.Load(indexName)
			if ok {
				disIndex.(*sync.Map).Delete(shardName)
			}
		}
		log.Debug().Str("type", meta.StorageEventTypeString[e.Type]).Str("key", string(e.Key)).Str("value", string(e.Value)).Msg("cluster distribution event")
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

func (c *Cluster) SetIndex(indexName string, metaVersion int64) error {
	c.indexes.Store(indexName, metaVersion)
	return metadata.Cluster.SetIndex(indexName, metaVersion)
}

func (c *Cluster) DeleteIndex(indexName string) error {
	// delete index from cluster
	c.indexes.Delete(indexName)
	_ = metadata.Cluster.DeleteIndex(indexName)

	// release distribution of this node
	disIndex, ok := c.distribution.Load(indexName)
	if ok {
		disIndex.(*sync.Map).Range(func(shard, value interface{}) bool {
			_ = metadata.Cluster.ReleaseDistribute(indexName, shard.(string))
			return true
		})
	}
	c.distribution.Delete(indexName)

	return nil
}

// DistributeIndexShards distribute shards by the index for the node id
func (c *Cluster) DistributeIndexShards(indexName string, nodeID int64) error {
	index, ok := GetIndex(indexName)
	if !ok {
		return errors.ErrIndexNotExists
	}
	needShards := int(math.Round(float64(index.GetShardNum()) / float64(c.GetNodesNum())))

	shards := make([]string, 0)
	disIndex, ok := c.distribution.Load(indexName)
	if !ok {
		disIndex = new(sync.Map)
		c.distribution.Store(indexName, disIndex)
		c.indexes.Store(indexName, index.GetMetaVersion())
	}
	disIndex.(*sync.Map).Range(func(shardName, value interface{}) bool {
		if nodeID == value.(int64) {
			shards = append(shards, shardName.(string))
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
			shards = append(shards, shardName.(string))
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
			disIndex.(*sync.Map).Store(shard, nid)
		} else {
			if oid, ok := disIndex.(*sync.Map).Load(shard); ok && oid == nodeID {
				continue
			}
			if len(shards) < needShards {
				if nid != nodeID {
					err := metadata.Cluster.ShardDistribute(indexName, shard, nodeID)
					if err != nil {
						return err
					}
				}
				// cache to local node
				disIndex.(*sync.Map).Store(shard, nodeID)
				// reload the shard from storage
				index.ReloadShard(shard)
				// add to shards
				index.localShards[shard] = index.shards[shard]
				index.shardHashing.Add(shard)
				shards = append(shards, shard)
			}
		}
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

	// cache to local node
	disIndex.(*sync.Map).Store(shard, nodeID)
	index.localShards[shard] = index.shards[shard]
	index.shardHashing.Add(shard)

	return nil
}

// ReleaseIndexShards release some shards let node just hold math.Round(shards / nodes)
func (c *Cluster) ReleaseIndexShards(indexName string, nodeID int64) error {
	index, ok := GetIndex(indexName)
	if !ok {
		return errors.ErrIndexNotExists
	}
	needShards := int(math.Round(float64(index.GetShardNum()) / float64(c.GetNodesNum())))

	disIndex, ok := c.distribution.Load(indexName)
	if !ok {
		return nil
	}
	shards := 0
	disIndex.(*sync.Map).Range(func(shardName, value interface{}) bool {
		if nodeID == value.(int64) {
			shards++
			if shards > needShards {
				// close shards
				index.shards[shardName.(string)].Close()
				// release distribution
				err := metadata.Cluster.ReleaseDistribute(indexName, shardName.(string))
				if err != nil {
					return false
				}
				delete(index.localShards, shardName.(string))
				index.shardHashing.Remove(shardName.(string))
				// clear cache
				disIndex.(*sync.Map).Delete(shardName)
			}
		}
		return true
	})

	return nil
}

// ReleaseNodeShards release some shards let node just hold math.Round(shards / nodes)
func (c *Cluster) ReleaseNodeShards(nodeID int64) {
	c.indexes.Range(func(indexName, version interface{}) bool {
		err := c.ReleaseIndexShards(indexName.(string), nodeID)
		if err != nil {
			log.Error().Err(err).Str("index", indexName.(string)).Msg("release shards distribution")
			return false
		}
		return true
	})
}
