package ider

import (
	"os"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/rs/zerolog/log"
)

var node *snowflake.Node

func init() {
	var err error
	nodeID := os.Getenv("ZINC_NODE_ID")
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
	return node.Generate().String()
}
