package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/ider"
)

func UpdateDocument(c *gin.Context) {
	indexName := c.Param("target")
	queryID := c.Param("id") // ID for the document to be updated provided in URL path

	var err error
	var doc map[string]interface{}
	if err = c.BindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	docID := ""
	mintedID := false

	// If id field is present then use it, else create a new UUID and use it
	if id, ok := doc["_id"]; ok {
		docID = id.(string)
	} else if queryID != "" {
		docID = queryID
	}
	if docID == "" {
		docID = ider.Generate()
		mintedID = true
	}

	// If the index does not exist, then create it
	index, exists := core.GetIndex(indexName)
	if !exists {
		index, err = core.NewIndex(indexName, "disk", core.UseNewIndexMeta, nil) // Create a new index with disk storage as default
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// store index
		core.StoreIndex(index)
	}

	err = index.UpdateDocument(docID, doc, mintedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": docID})
}
