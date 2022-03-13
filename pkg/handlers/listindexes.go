package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/core"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func ListIndexes(c *gin.Context) {
	items := make([]*meta.Index, 0, len(core.ZINC_INDEX_LIST))
	for name, value := range core.ZINC_INDEX_LIST {
		item := new(meta.Index)
		item.Name = name
		item.StorageType = value.StorageType
		if value.Settings != nil {
			item.Settings = value.Settings
		} else {
			item.Settings = new(meta.IndexSettings)
		}
		if value.Mappings != nil {
			item.Mappings = value.CachedMappings
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)
}
