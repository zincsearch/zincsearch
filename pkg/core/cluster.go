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
	"errors"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"

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

func (c *Cluster) GetNodesLen() int {
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
			// need release some shards
			if err := c.ReleaseShardsDistribution(c.Local().ID); err != nil {
				log.Error().Err(err).Int64("nodeid", c.Local().ID).Msg("release shards distribution")
			}
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
			version, _ := strconv.ParseInt(string(e.Value), 10, 64)
			oldVersion, ok := c.indexes.Load(indexName)
			if ok && oldVersion.(int64) == version {
				continue // no change
			}
			// need reload index
			ZINC_INDEX_LIST.Delete(indexName)
			err := LoadIndexFromMetadata(indexName, meta.Version)
			if err != nil {
				log.Error().Str("index", indexName).Int64("version", version).Msg("cluster index event, add index")
				continue
			}
			c.indexes.Store(indexName, version)
		case meta.StorageEventTypeDelete:
			err := c.DeleteIndex(indexName)
			if err != nil {
				log.Error().Str("index", indexName).Int64("nodeid", c.Local().ID).Msg("cluster index event, delete index")
			}
			_ = DeleteIndex(indexName)
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
	lock, err := metadata.Cluster.NewLocker("meta/nodes")
	if err != nil {
		return err
	}
	lock.Lock()
	defer lock.Unlock()

	// get nodes and calculate the new node id
	nodes, err := metadata.Cluster.ListNodes(0, 0)
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
	ZINC_NODE.SetID(newNodeID)

	// cache local node
	c.nodes.Store(newNodeID, ZINC_NODE.data)

	// add node
	return metadata.Cluster.Join(newNodeID, c.Local())
}

func (c *Cluster) Leave() error {
	return metadata.Cluster.Leave(c.Local().ID)
}

func (c *Cluster) StoreIndex(indexName string, metaVersion int64) {
	c.indexes.Store(indexName, metaVersion)
}

func (c *Cluster) DeleteIndex(indexName string) error {
	c.indexes.Delete(indexName)
	c.distribution.Delete(indexName)
	return metadata.Cluster.DeleteIndex(indexName)
}

// GetShardsDistribution return the shards distribution of the index by the node id
func (c *Cluster) GetShardsDistribution(indexName string, nodeID int64) ([]string, error) {
	index, ok := GetIndex(indexName)
	if !ok {
		return nil, errors.New("index not found")
	}
	needShards := int(math.Round(float64(index.shardNum) / float64(c.GetNodesLen())))

	shards := make([]string, 0)
	disIndex, ok := c.distribution.Load(indexName)
	if !ok {
		disIndex = new(sync.Map)
		c.distribution.Store(indexName, disIndex)
		c.indexes.Store(indexName, index.ref.MetaVersion)
	}
	disIndex.(*sync.Map).Range(func(shardName, value interface{}) bool {
		if nodeID == value.(int64) {
			shards = append(shards, shardName.(string))
		}
		return true
	})
	if len(shards) >= needShards {
		return shards, nil
	}

	// no local cache, fetch from metadata
	lock, err := metadata.Cluster.NewLocker("meta/distribution/" + indexName)
	if err != nil {
		return nil, err
	}
	lock.Lock()
	defer lock.Unlock()
	metaShards, err := metadata.Cluster.ListShards(indexName)
	if err != nil {
		return nil, err
	}

	for id := range index.shards {
		if _, ok := metaShards[id]; !ok {
			metaShards[id] = 0
		}
	}

	for shard, nid := range metaShards {
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
						return nil, err
					}
				}
				disIndex.(*sync.Map).Store(shard, nodeID)
				// add to shards
				index.localShards[shard] = index.shards[shard]
				index.shardHashing.Add(shard)
				shards = append(shards, shard)
			}
		}
	}

	return shards, nil
}

// ReleaseShardsDistribution release some shards let node just hold math.Round(shards / nodes)
func (c *Cluster) ReleaseShardsDistribution(nodeID int64) error {
	c.indexes.Range(func(indexName, version interface{}) bool {
		err := c.ReleaseShardsDistributionByIndex(indexName.(string), nodeID)
		if err != nil {
			log.Error().Err(err).Str("index", indexName.(string)).Int64("nodeid", nodeID).Msg("release shards distribution")
			return false
		}
		return true
	})
	return nil
}

// ReleaseShardsDistributionByIndex release some shards let node just hold math.Round(shards / nodes)
func (c *Cluster) ReleaseShardsDistributionByIndex(indexName string, nodeID int64) error {
	index, ok := GetIndex(indexName)
	if !ok {
		return errors.New("index not found")
	}
	needShards := int(math.Round(float64(index.shardNum) / float64(c.GetNodesLen())))

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
