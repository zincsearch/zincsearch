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
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sony/sonyflake/v2"

	"github.com/zincsearch/zincsearch/pkg/config"
	"github.com/zincsearch/zincsearch/pkg/zutils/base62"
)

// epoch is the fixed start time of the ID clock; changing it changes the generated IDs.
var epoch = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

type Node struct {
	node *sonyflake.Sonyflake
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

// maxNodeID is the largest node id representable in 10 machine bits.
const maxNodeID = 1<<10 - 1

// NewNode returns a generator laid out as 41 bits of milliseconds, 12 bits of sequence, 10 bits of node (low bits).
// id must be in [0, maxNodeID]; out-of-range ids are rejected rather than folded onto another node.
func NewNode(id int) (*Node, error) {
	if id < 0 || id > maxNodeID {
		return nil, fmt.Errorf("node id %d out of range [0, %d]", id, maxNodeID)
	}
	node, err := sonyflake.New(sonyflake.Settings{
		BitsSequence:  12,
		BitsMachineID: 10,
		TimeUnit:      time.Millisecond,
		StartTime:     epoch,
		MachineID:     func() (int, error) { return id, nil },
	})
	if err != nil {
		return nil, err
	}
	return &Node{node: node}, nil
}

func (n *Node) Generate() string {
	id, err := n.node.NextID()
	if err != nil {
		log.Fatal().Err(err).Msg("id generate failed")
	}
	return base62.Encode(id)
}
