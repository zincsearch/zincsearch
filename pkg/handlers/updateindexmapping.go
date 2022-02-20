package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/core"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/mappings"
)

func UpdateIndexMapping(c *gin.Context) {
	indexName := c.Param("target")
	if indexName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index.name should be not empty"})
		return
	}

	_, exists := core.GetIndex(indexName)
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index [" + indexName + "] already exists"})
		return
	}

	var newIndex core.Index
	c.BindJSON(&newIndex)
	mappings, err := mappings.Request(newIndex.Mappings)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	index, err := core.NewIndex(indexName, newIndex.StorageType, core.UseNewIndexMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// update mappings
	if mappings != nil && len(mappings.Properties) > 0 {
		index.SetMappings(mappings)
	}

	// store index
	core.StoreIndex(index)

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
