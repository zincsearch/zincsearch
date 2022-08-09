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

package metadata

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/goccy/go-json"

	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/metadata/storage"
)

type cluster struct{}

var Cluster = new(cluster)

func (t *cluster) NewLocker(key string) (sync.Locker, error) {
	return db.NewLocker(t.key("lock/" + key))
}

func (t *cluster) ListNode(offset, limit int64) ([]*meta.Node, error) {
	data, err := db.List(t.key("node/"), offset, limit)
	if err != nil {
		return nil, err
	}
	nodes := make([]*meta.Node, 0, len(data))
	for _, d := range data {
		node := new(meta.Node)
		err = json.Unmarshal(d, node)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })

	return nodes, nil
}

func (t *cluster) ListIndex(offset, limit int64) (map[string]int64, error) {
	data, err := db.ListEntries(t.key("index/"), offset, limit)
	if err != nil {
		return nil, err
	}
	indexes := make(map[string]int64, len(data))
	for _, d := range data {
		metaVersion, _ := strconv.ParseInt(string(d.Value), 10, 64)
		if err != nil {
			return nil, err
		}
		indexes[string(d.Key)] = metaVersion
	}
	return indexes, nil
}

// ListDistribution returns the distribution of shards for an index.
// returns data is map[shardName]nodeID
func (t *cluster) ListDistribution(index string) (map[string]int64, error) {
	data, err := db.ListEntries(t.key("distribution/"+index+"/"), 0, 0)
	if err != nil {
		return nil, err
	}
	shards := make(map[string]int64, len(data))
	for _, d := range data {
		shard := string(d.Key)
		nodeID, _ := strconv.ParseInt(string(d.Value), 10, 64)
		shards[shard] = nodeID
	}
	return shards, nil
}

func (t *cluster) Join(id int64, node *meta.Node) error {
	data, err := json.Marshal(node)
	if err != nil {
		return err
	}
	return db.SetWithKeepAlive(t.key(fmt.Sprintf("node/%d", id)), data, 5)
}

func (t *cluster) Leave(id int64) error {
	return db.Delete(t.key(fmt.Sprintf("node/%d", id)))
}

func (t *cluster) ShardDistribute(index, shard string, nodeID int64) error {
	data := strconv.FormatInt(nodeID, 10)
	return db.SetWithKeepAlive(t.key(fmt.Sprintf("distribution/%s/%s", index, shard)), []byte(data), 5)
}

func (t *cluster) ReleaseDistribute(index, shard string) error {
	return db.CancelWithKeepAlive(t.key(fmt.Sprintf("distribution/%s/%s", index, shard)))
}

func (t *cluster) GetShardDistribute(index, shard string) (int64, error) {
	data, err := db.Get(t.key(fmt.Sprintf("distribution/%s/%s", index, shard)))
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(data), 10, 64)
}

func (t *cluster) SetIndex(index string, metaVersion int64) error {
	data := strconv.FormatInt(metaVersion, 10)
	return db.Set(t.key("index/"+index), []byte(data))
}

func (t *cluster) GetIndex(index string) (int64, error) {
	data, err := db.Get(t.key("index/" + index))
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(data), 10, 64)
}

func (t *cluster) DeleteIndex(index string) error {
	return db.Delete(t.key("index/" + index))
}

func (t *cluster) Watch(eventType string) <-chan storage.StorageEvent {
	key := ""
	switch eventType {
	case meta.ClusterEventTypeNode:
		key = "node/"
	case meta.ClusterEventTypeIndex:
		key = "index/"
	case meta.ClusterEventTypeDistribution:
		key = "distribution/"
	}
	return db.Watch(t.key(key))
}

func (t *cluster) key(key string) string {
	return "/cluster/" + key
}
