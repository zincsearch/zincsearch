package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/core"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

func SearchIndex(c *gin.Context) {

	indexName := c.Param("target")

	// fmt.Println("Got search request for index: ", indexName)

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
