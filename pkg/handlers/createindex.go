package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/core"
)

func CreateIndex(c *gin.Context) {
	var newIndex core.Index
	c.BindJSON(&newIndex)

	if newIndex.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "index.name should be not empty"})
		return
	}

	if ok, storage := core.IndexExists(newIndex.Name); ok {
		c.JSON(http.StatusOK, gin.H{
			"result":       "Index: " + newIndex.Name + " exists",
			"storage_type": storage,
		})
		return
	}

	cIndex, err := core.NewIndex(newIndex.Name, newIndex.StorageType)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	core.ZINC_INDEX_LIST[newIndex.Name] = cIndex

	c.JSON(http.StatusOK, gin.H{
		"result":       "Index: " + newIndex.Name + " created",
		"storage_type": newIndex.StorageType,
	})
}
