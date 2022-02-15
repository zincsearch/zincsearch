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

	mappings, err := mappings.Request(newIndex.Mappings)
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

		// use template
		if mappings == nil {
			template, _ := core.UseTemplate(newIndex.Name)
			if template != nil && template.Template.Mappings != nil && len(template.Template.Mappings.Properties) > 0 {
				mappings = template.Template.Mappings
			}
		}
	}

	// update mappings
	if mappings != nil && len(mappings.Properties) > 0 {
		index.SetMappings(mappings)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "index " + newIndex.Name + " created",
		"storage_type": newIndex.StorageType,
	})
}
