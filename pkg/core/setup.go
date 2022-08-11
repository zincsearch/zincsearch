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

import "github.com/zinclabs/zinc/pkg/meta"

// init initializes the core package.
func init() {
	SetupNode()
	SetupCluster()
	SetupShardWAL()
	SetupIndex()
	// load distribution of cluster
	ZINC_CLUSTER.LoadDistribution()
	// ready to serve requests
	ZINC_NODE.SetStatus(meta.NodeStatusOK)
}
