package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/core"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/mappings"
)

func CreateIndex(c *gin.Context) {
	var newIndex core.Index
	c.BindJSON(&newIndex)

	if newIndex.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index.name should be not empty"})
		return
	}

	if _, ok := core.GetIndex(newIndex.Name); ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index [" + newIndex.Name + "] already exists"})
		return
	}

	mappings, err := mappings.Request(newIndex.Mappings)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	index, err := core.NewIndex(newIndex.Name, newIndex.StorageType, core.UseNewIndexMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// update settings
	index.SetSettings(newIndex.Settings)

	// update mappings
	index.SetMappings(mappings)

	// store index
	err = core.StoreIndex(index, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// store index
	core.StoreIndex(index, false)

	c.JSON(http.StatusOK, gin.H{
		"message":      "index created",
		"index":        newIndex.Name,
		"storage_type": newIndex.StorageType,
	})
}
