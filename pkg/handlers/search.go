package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/core"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

// SearchIndex searches the index for the given http request from end user
func SearchIndex(c *gin.Context) {
	indexName := c.Param("target")
	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index " + indexName + " does not exists"})
		return
	}

	var iQuery v1.ZincQuery
	if err := c.BindJSON(&iQuery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := index.Search(&iQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event_data := make(map[string]interface{})
	event_data["search_type"] = iQuery.SearchType
	event_data["search_index_storage"] = index.StorageType
	event_data["search_index_size_in_mb"] = index.Size
	event_data["time_taken_to_search_in_ms"] = res.Took
	event_data["aggregations_count"] = len(iQuery.Aggregations)
	core.Telemetry.Event("search", event_data)

	c.JSON(http.StatusOK, res)
}
