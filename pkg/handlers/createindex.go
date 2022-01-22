package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/core"
)

func CreateIndex(c *gin.Context) {
	var newIndex core.Index
	c.BindJSON(&newIndex)

	if newIndex.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "index.name should be not empty"})
		return
	}

	mapping := make(map[string]string)
	for field, prop := range newIndex.Mappings.Properties {
		ptype := strings.ToLower(prop.Type)
		switch ptype {
		case "text", "keyword", "numeric", "bool", "time":
			ptype = ptype
		case "boolean":
			ptype = "bool"
		case "date", "datetime":
			ptype = "time"
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "mappings unsupport type: " + prop.Type})
			return
		}
		mapping[field] = ptype
	}

	var cIndex *core.Index
	var ok bool
	var err error
	if cIndex, ok = core.IndexExists(newIndex.Name); !ok {
		cIndex, err = core.NewIndex(newIndex.Name, newIndex.StorageType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		core.ZINC_INDEX_LIST[newIndex.Name] = cIndex
	}

	// update mapping
	if len(mapping) > 0 {
		cIndex.SetMapping(mapping)
	}

	c.JSON(http.StatusOK, gin.H{
		"result":       "index " + newIndex.Name + " created",
		"storage_type": newIndex.StorageType,
	})
}
