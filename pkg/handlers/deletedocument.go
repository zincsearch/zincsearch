package handlers

import (
	"net/http"

	"github.com/blugelabs/bluge"
	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/core"
)

func DeleteDocument(c *gin.Context) {
	indexName := c.Param("target")
	queryID := c.Param("id")

	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index does not exists"})
		return
	}

	bdoc := bluge.NewDocument(queryID)
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil))
	err := index.Writer.Delete(bdoc.ID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		core.ZINC_INDEX_LIST[indexName].ReduceDocsCount(1)
		c.JSON(http.StatusOK, gin.H{"message": "Deleted", "index": indexName, "id": queryID})
	}
}
