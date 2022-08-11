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

package meta

// Default field name
const (
	TimeFieldName   = "@timestamp"
	IDFieldName     = "@_id"
	ActionFieldName = "@_action"
	ShardFieldName  = "@_shard"
)

const (
	ActionTypeInsert = "insert"
	ActionTypeUpdate = "update"
	ActionTypeDelete = "delete"
)

const (
	ServerModeNode    = "node"
	ServerModeCluster = "cluster"
)

const (
	StorageEventTypePut    = int64(1)
	StorageEventTypeDelete = int64(2)
	StorageEventTypeCreate = int64(3)
	StorageEventTypeUpdate = int64(4)
)

var StorageEventTypeString = map[int64]string{
	StorageEventTypePut:    "PUT",
	StorageEventTypeDelete: "DELETE",
	StorageEventTypeCreate: "CREATE",
	StorageEventTypeUpdate: "UPDATE",
}

const (
	ClusterEventTypeNode         = "node"         // node event
	ClusterEventTypeIndex        = "index"        // index event
	ClusterEventTypeDistribution = "distribution" // distribute event

	ClusterEventNodeJoin  = "join"  // node join cluster
	ClusterEventNodeLeave = "leave" // node leave cluster

	ClusterEventIndexCreate = "create" // index create event
	ClusterEventIndexRemove = "remove" // index remove event
	ClusterEventIndexUpdate = "update" // index metadata update event

	ClusterEventDistributionHold = "hold" // node hold a shard event
	ClusterEventDistributionFree = "free" // node free a shard event
)
