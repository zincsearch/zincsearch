package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/core"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

// SearchIndex searches the index for the given http request from end user
func SearchIndex(c *gin.Context) {

	indexName := c.Param("target")
	indexExists, _ := core.IndexExists(indexName)

	if !indexExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index '" + indexName + "' does not exist"})
		return
	}

	var iQuery v1.ZincQuery

	c.BindJSON(&iQuery)

	index := core.ZINC_INDEX_LIST[indexName]
	res, errS := index.Search(iQuery)

	if errS != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errS.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}
