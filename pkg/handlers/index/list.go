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
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zinc/pkg/core"
)

// @Id ListIndexes
// @Summary List indexes
// @Tags    Index
// @Produce json
// @Success 200 {object} []core.Index
// @Router /api/index [get]
func List(c *gin.Context) {
	items := core.ZINC_INDEX_LIST.ListStat()
	c.JSON(http.StatusOK, items)
}

// @Id IndexNameList
// @Summary List index Name
// @Tags    Index
// @Param   name  query  string  false  "IndexName"
// @Produce json
// @Success 200 {object} []string
// @Router /api/index_name [get]
func IndexNameList(c *gin.Context) {
	queryName := strings.ToLower(c.DefaultQuery("name", ""))
	var items []string
	names := core.ZINC_INDEX_LIST.ListName()
	if queryName == "" {
		items = names
	} else {
		for _, name := range names {
			if strings.Contains(strings.ToLower(name), queryName) {
				items = append(items, name)
			}
		}
	}

	count := 30
	if len(items) > count {
		c.JSON(http.StatusOK, items[0:count])
	} else {
		c.JSON(http.StatusOK, items)
	}
}
