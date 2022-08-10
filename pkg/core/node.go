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
	"os"
	"sync/atomic"

	"github.com/zinclabs/zinc/pkg/meta"
)

var ZINC_NODE *Node

type Node struct {
	data *meta.Node
}

func SetupNode() {
	name, _ := os.Hostname()
	ZINC_NODE = &Node{
		data: &meta.Node{
			ID:     1,
			Name:   name,
			Status: meta.NodeStatusPrepare,
		},
	}
}

func (n *Node) SetID(id int64) {
	atomic.StoreInt64(&n.data.ID, id)
}

func (n *Node) GetID() int64 {
	return atomic.LoadInt64(&n.data.ID)
}

func (n *Node) GetStatus() string {
	status := atomic.LoadInt64(&n.data.Status)
	return meta.NodeStatusString[status]
}

func (n *Node) SetStatus(status int64) {
	atomic.StoreInt64(&n.data.Status, status)
}
