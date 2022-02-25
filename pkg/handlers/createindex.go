package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/core"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	zincanalysis "github.com/prabhatsharma/zinc/pkg/uquery/v2/analysis"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/mappings"
)

func CreateIndex(c *gin.Context) {
	var newIndex core.Index
	if err := c.BindJSON(&newIndex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	indexName := c.Param("target")
	if newIndex.Name == "" && indexName != "" {
		newIndex.Name = indexName
	}

	if newIndex.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index.name should be not empty"})
		return
	}

	if _, ok := core.GetIndex(newIndex.Name); ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index [" + newIndex.Name + "] already exists"})
		return
	}

	if newIndex.Settings == nil {
		newIndex.Settings = meta.NewIndexSettings()
	}
	analyzers, err := zincanalysis.RequestAnalyzer(newIndex.Settings.Analysis)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// update analyzers
	index.SetAnalyzers(analyzers)

	// update mappings
	index.SetMappings(mappings)

	// store index
	err = core.StoreIndex(index)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "index created",
		"index":        newIndex.Name,
		"storage_type": newIndex.StorageType,
	})
}
