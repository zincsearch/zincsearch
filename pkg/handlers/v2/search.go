package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc/pkg/core"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

// SearchIndex searches the index for the given http request from end user
func SearchIndex(c *gin.Context) {
	indexName := c.Param("target")
	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index " + indexName + " does not exists"})
		return
	}

	query := new(meta.ZincQuery)
	err := c.BindJSON(query)
	if err != nil {
		log.Printf("handlers.SearchIndex: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	resp, err := index.SearchV2(query)
	if resp != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if err != nil {
		switch v := err.(type) {
		case *meta.Error:
			c.JSON(http.StatusBadRequest, v)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": v})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}
