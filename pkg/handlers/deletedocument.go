package handlers

import (
	"net/http"

	"github.com/blugelabs/bluge"
	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/core"
)

func DeleteDocument(c *gin.Context) {
	indexName := c.Param("target")
	query_id := c.Param("id")

	indexExists, _ := core.IndexExists(indexName)

	if !indexExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index not exist"})
		return
	}

	// log.Printf("deleet document indexName:%[1]s, query_id:%[2]s", indexName, query_id)

	bdoc := bluge.NewDocument(query_id)

	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil))

	docIndexWriter := core.ZINC_INDEX_LIST[indexName].Writer

	err := docIndexWriter.Delete(bdoc.ID())

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)

	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Deleted", "index": indexName, "id": query_id})
	}
}
