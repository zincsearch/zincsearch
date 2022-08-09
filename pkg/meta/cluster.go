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

// Cluster is a collection of nodes.
// -- Nodes contains nodes in the cluster.
// -- Distribution contains [index:shards] of nodes in the cluster.
type Cluster struct {
	Nodes        map[int64]*Node             `json:"nodes"`        // nodeID -> node
	Indexes      map[string]int64            `json:"indexes"`      // indexName -> index metadata version
	Distribution map[string]map[string]int64 `json:"distribution"` // index -> indexShard (first layer) -> nodeID

}

func NewCluster() *Cluster {
	return &Cluster{
		Nodes:        make(map[int64]*Node, 3),
		Indexes:      make(map[string]int64, 2),
		Distribution: make(map[string]map[string]int64, 2),
	}
}
