package handlers

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/core"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func ListIndexes(c *gin.Context) {
	items := make(meta.SortIndex, 0, len(core.ZINC_INDEX_LIST))
	for name, value := range core.ZINC_INDEX_LIST {
		item := new(meta.Index)
		item.Name = name
		item.StorageType = value.StorageType
		item.StorageSize = int64(value.StorageSize)
		item.DocsCount = value.DocsCount
		if value.Settings != nil {
			item.Settings = value.Settings
		} else {
			item.Settings = new(meta.IndexSettings)
		}
		if value.CachedMappings != nil {
			// format mappings
			mappings := value.CachedMappings
			if mappings == nil {
				mappings = meta.NewMappings()
			}
			item.Mappings = mappings
		}
		items = append(items, item)
	}

	sort.Sort(items)

	c.JSON(http.StatusOK, items)
}
