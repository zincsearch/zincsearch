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

	var newIndex core.Index
	c.BindJSON(&newIndex)
	mappings, err := mappings.Request(newIndex.Mappings)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

		// use template
		template, _ := core.UseTemplate(indexName)
		if template != nil && template.Template.Mappings != nil {
			mappings = template.Template.Mappings
		}
	}

	// update mappings
	if mappings != nil && len(mappings.Properties) > 0 {
		index.SetMappings(mappings)
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
