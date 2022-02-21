package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/core"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/analyzer"
)

func Analyze(c *gin.Context) {
	var ana analyzeRequest
	if err := c.BindJSON(&ana); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	indexName := c.Param("target")
	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "index " + indexName + " does not exists"})
		return
	}

	analyzer, err := analyzer.Query(index.CachedAnalysis, ana.Analyzer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "analyzer " + ana.Analyzer + " does not exists"})
		return
	}

	tokens := analyzer.Analyze([]byte(ana.Text))
	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}

type analyzeRequest struct {
	Analyzer string `json:"analyzer"`
	Text     string `json:"text"`
}
