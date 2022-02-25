package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/core"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	zincanalysis "github.com/prabhatsharma/zinc/pkg/uquery/v2/analysis"
)

func GetIndexSettings(c *gin.Context) {
	indexName := c.Param("target")
	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "index " + indexName + " does not exists"})
		return
	}

	settings := index.Settings
	if settings == nil {
		settings = new(meta.IndexSettings)
	}

	c.JSON(http.StatusOK, gin.H{index.Name: gin.H{"settings": settings}})
}

func UpdateIndexSettings(c *gin.Context) {
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

	if newIndex.Settings == nil {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
		return
	}

	analyzers, err := zincanalysis.RequestAnalyzer(newIndex.Settings.Analysis)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	index, exists := core.GetIndex(indexName)
	if exists {
		// it can only change settings.NumberOfReplicas when index exists
		if newIndex.Settings.NumberOfReplicas > 0 {
			index.Settings.NumberOfReplicas = newIndex.Settings.NumberOfReplicas
		}
		if newIndex.Settings.Analysis != nil && len(newIndex.Settings.Analysis.Analyzer) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't update analyzer for existing index"})
			return
		}
		// store index
		core.StoreIndex(index)

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
		return
	}

	index, err = core.NewIndex(indexName, newIndex.StorageType, core.UseNewIndexMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// update settings
	index.SetSettings(newIndex.Settings)

	// update analyzers
	index.SetAnalyzers(analyzers)

	// store index
	core.StoreIndex(index)

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
