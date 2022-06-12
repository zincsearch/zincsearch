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

package index

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/meta"
)

// @Summary List Indexes
// @Tags  Index
// @Produce json
// @Success 200 {object} []core.Index
// @Router /api/index [get]
func List(c *gin.Context) {
	items := core.ZINC_INDEX_LIST.List()
	for _, index := range items {
		if index.Settings == nil {
			index.Settings = new(meta.IndexSettings)
		}
		if index.Mappings == nil {
			index.Mappings = meta.NewMappings()
		}
		// update metadata while listing
		index.UpdateMetadata()
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	c.JSON(http.StatusOK, items)
}
