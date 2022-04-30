/* Copyright 2022 Zinc Labs Inc. and Contributors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package handlers

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/zinclabs/zinc/pkg/core"
	meta "github.com/zinclabs/zinc/pkg/meta/v2"
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
