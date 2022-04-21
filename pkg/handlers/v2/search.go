package v2

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/errors"
	meta "github.com/zinclabs/zinc/pkg/meta/v2"
)

// SearchIndex searches the index for the given http request from end user
func SearchIndex(c *gin.Context) {
	indexName := c.Param("target")

	query := new(meta.ZincQuery)
	if err := c.BindJSON(query); err != nil {
		log.Printf("handlers.v2.SearchIndex: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storageType := "disk"
	indexSize := 0.0

	var err error
	var resp *meta.SearchResponse
	if indexName == "" || strings.HasSuffix(indexName, "*") {
		resp, err = core.MultiSearchV2(indexName, query)
	} else {
		index, exists := core.GetIndex(indexName)
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "index " + indexName + " does not exists"})
			return
		}

		storageType = index.StorageType
		indexSize = index.StorageSize
		resp, err = index.SearchV2(query)
	}

	if err != nil {
		handleError(c, err)
		return
	}

	eventData := make(map[string]interface{})
	eventData["search_type"] = "query_dsl"
	eventData["search_index_storage"] = storageType
	eventData["search_index_size_in_mb"] = indexSize
	eventData["time_taken_to_search_in_ms"] = resp.Took
	eventData["aggregations_count"] = len(query.Aggregations)
	core.Telemetry.Event("search", eventData)

	c.JSON(http.StatusOK, resp)
}

func handleError(c *gin.Context, err error) {
	if err != nil {
		switch v := err.(type) {
		case *errors.Error:
			c.JSON(http.StatusBadRequest, gin.H{"error": v})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": v.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
