package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/core"
)

func ListIndexes(c *gin.Context) {
	var indexListMap = make(map[string]*SimpleIndex)
	for name, value := range core.ZINC_INDEX_LIST {
		mappings := make(map[string]string, len(value.CachedMappings.Properties))
		for field, prop := range value.CachedMappings.Properties {
			if field == "_id" || field == "@timestamp" {
				continue
			}
			mappings[field] = prop.Type
		}
		indexListMap[name] = &SimpleIndex{
			Name:     name,
			Mappings: mappings,
		}
	}
	c.JSON(http.StatusOK, indexListMap)
}

type SimpleIndex struct {
	Name     string            `json:"name"`
	Mappings map[string]string `json:"mapping"`
}
