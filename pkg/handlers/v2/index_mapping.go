package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/core"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/mappings"
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
		mappings = meta.NewMappings()
	}
	for field := range mappings.Properties {
		if field == "_id" || field == "@timestamp" {
			delete(mappings.Properties, field)
		}
	}

	c.JSON(http.StatusOK, gin.H{index.Name: gin.H{"mappings": mappings}})
}

func UpdateIndexMapping(c *gin.Context) {
	indexName := c.Param("target")
	if indexName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index.name should be not empty"})
		return
	}

	var newIndex core.Index
	if err := c.BindJSON(&newIndex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, exists := core.GetIndex(indexName)
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index [" + indexName + "] already exists"})
		return
	}

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
