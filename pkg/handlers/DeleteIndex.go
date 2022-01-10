package handlers

import (
	"net/http"
	"os"

	"github.com/blugelabs/bluge"
	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/core"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

// DeleteIndex deletes a zinc index and its associated data. Be careful using thus as you ca't undo this action.
func DeleteIndex(c *gin.Context) {
	indexName := c.Param("indexName")

	if _, ok := core.ZINC_INDEX_LIST[indexName]; !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "index not exist"})
		return
	}

	// 1. Close the index writer
	core.ZINC_INDEX_LIST[indexName].Writer.Close()

	// 2. Delete from the cache
	delete(core.ZINC_INDEX_LIST, indexName)

	// 3. Physically delete the index
	DATA_PATH := zutils.GetEnv("DATA_PATH", "./data")

	err := os.RemoveAll(DATA_PATH + "/" + indexName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 4. Delete the index mapping
	bdoc := bluge.NewDocument(indexName)
	err = core.ZINC_SYSTEM_INDEX_LIST["_index_mapping"].Writer.Delete(bdoc.ID())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Deleted",
			"index":   indexName,
		})
	}
}
