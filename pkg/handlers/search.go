package handlers

import (
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc/pkg/core"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

// SearchIndex searches the index for the given http request from end user
func SearchIndex(c *gin.Context) {
	start_time := time.Now()

	indexName := c.Param("target")
	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index " + indexName + " does not exists"})
		return
	}

	var iQuery v1.ZincQuery
	err := c.BindJSON(&iQuery)
	if err != nil {
		log.Printf("handlers.SearchIndex: %v", err)
		return
	}

	res, err := index.Search(&iQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	total_time_taken := time.Since(start_time)

	event_data := make(map[string]interface{})
	event_data["search_type"] = iQuery.SearchType
	event_data["search_index_storage"] = core.ZINC_INDEX_LIST[indexName].IndexType
	event_data["search_index_size_in_mb"] = math.Round(core.GetIndexSize(indexName))
	event_data["time_taken_to_search_in_ms"] = total_time_taken / 1000 / 1000
	event_data["aggregations_count"] = len(iQuery.Aggregations)
	core.TelemetryEvent("search", event_data)

	c.JSON(http.StatusOK, res)
}
