package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/core"
)

func GetIndexMappings(c *gin.Context) {
	indexName := c.Param("target")
	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index " + indexName + " does not exists"})
		return
	}

	mappings, err := index.GetStoredMappings()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// format mappings
	for field := range mappings.Properties {
		if field == "_id" || field == "@timestamp" {
			delete(mappings.Properties, field)
		}
	}

	c.JSON(http.StatusOK, mappings)
}
