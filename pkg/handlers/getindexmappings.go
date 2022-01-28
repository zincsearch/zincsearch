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

	mappings, err := index.GetStoredMapping()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// format mapping
	properties := make(map[string]core.Properties, len(mappings))
	for field, pType := range mappings {
		if field == "_id" || field == "@timestamp" {
			continue
		}
		properties[field] = core.Properties{Type: pType}
	}

	c.JSON(http.StatusOK, gin.H{"mappings": gin.H{"properties": properties}})
}
