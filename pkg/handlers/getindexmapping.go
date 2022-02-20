package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/core"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func GetIndexMapping(c *gin.Context) {
	indexName := c.Param("target")
	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "index " + indexName + " does not exists"})
		return
	}

	// format mappings
	mappings := index.CachedMappings
	if mappings == nil {
		mappings = new(meta.Mappings)
	} else {
		for field := range mappings.Properties {
			if field == "_id" || field == "@timestamp" {
				delete(mappings.Properties, field)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{index.Name: gin.H{"mappings": mappings}})
}
