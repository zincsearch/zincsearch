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
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/zutils/base62"
)

var node *snowflake.Node

func init() {
	var err error
	nodeID := config.Global.NodeID
	if nodeID == "" {
		nodeID = "1"
	}
	id, _ := strconv.ParseInt(nodeID, 10, 64)
	node, err = snowflake.NewNode(id)
	if err != nil {
		log.Fatal().Msgf("id generater init err %s", err.Error())
	}
}

func Generate() string {
	return base62.Encode(node.Generate().Int64())
}
