package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/core"
)

func CreateIndex(c *gin.Context) {
	fmt.Println("CreateIndex")

	var newIndex core.Index
	c.BindJSON(&newIndex)

	cIndex, err := core.NewIndex(newIndex.Name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	core.ZINC_INDEX_LIST[newIndex.Name] = cIndex

	c.JSON(http.StatusOK, gin.H{
		"result": "Index: " + newIndex.Name + " created",
	})
}
