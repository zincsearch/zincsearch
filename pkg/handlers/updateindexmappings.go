package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/core"
)

func UpdateIndexMappings(c *gin.Context) {
	indexName := c.Param("target")
	if indexName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "index.name should be not empty"})
		return
	}

	var newIndex core.Index
	c.BindJSON(&newIndex)
	mappings, err := core.FormatMapping(&newIndex.Mappings)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	index, exists := core.GetIndex(indexName)
	if !exists {
		index, err = core.NewIndex(indexName, newIndex.StorageType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		core.ZINC_INDEX_LIST[indexName] = index
	}

	// update mapping
	if len(mappings) > 0 {
		index.SetMapping(mappings)
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
