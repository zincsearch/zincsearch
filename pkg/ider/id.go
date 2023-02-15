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

package ider

import (
	"github.com/bwmarrin/snowflake"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zincsearch/pkg/config"
	"github.com/zinclabs/zincsearch/pkg/zutils/base62"
)

type Node struct {
	node *snowflake.Node
}

var local *Node

func init() {
	var err error
	local, err = NewNode(config.Global.NodeID)
	if err != nil {
		log.Fatal().Msgf("id generater init[local] err %s", err.Error())
	}
}

func Generate() string {
	return local.Generate()
}

func NewNode(id int) (*Node, error) {
	node, err := snowflake.NewNode(int64(id % 1024))
	return &Node{node: node}, err
}

func (n *Node) Generate() string {
	return base62.Encode(n.node.Generate().Int64())
}
