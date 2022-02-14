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
		c.JSON(http.StatusBadRequest, gin.H{"error": "index.name should be not empty"})
		return
	}

	mappings, err := core.FormatMappings(*newIndex.Mappings)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ok bool
	var index *core.Index
	if index, ok = core.GetIndex(newIndex.Name); !ok {
		index, err = core.NewIndex(newIndex.Name, newIndex.StorageType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		core.ZINC_INDEX_LIST[newIndex.Name] = index
	}

	// update mappings
	if len(mappings.Properties) > 0 {
		index.SetMappings(mappings)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "index " + newIndex.Name + " created",
		"storage_type": newIndex.StorageType,
	})
}
