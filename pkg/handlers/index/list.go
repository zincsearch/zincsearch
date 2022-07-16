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
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/meta"
)

// @Id ListIndexes
// @Summary List indexes
// @Tags    Index
// @Param   page_num  query  integer  false  "page num"
// @Param   page_size  query  integer  false  "page size"
// @Param   sort_by  query  string  false  "sort by"
// @Param   descending  query  boolen  false  "descending"
// @Param   filter  query  string  false  "filter"
// @Produce json
// @Success 200 {object} {list:[]core.Index, page: meta.Page}
// @Router /api/index [get]
func List(c *gin.Context) {
	page := meta.NewPage(c)
	sortBy := c.DefaultQuery("sort_by", "")
	descending, _ := strconv.ParseBool(c.DefaultQuery("descending", "false"))
	filter := c.DefaultQuery("filter", "")

	items := core.ZINC_INDEX_LIST.ListStat()

	if len(filter) > 0 {
		var res []*core.Index
		for _, item := range items {
			if strings.Contains(item.Name, filter) {
				res = append(res, item)
			}
		}
		items = res
	}

	if len(sortBy) > 0 {
		if sortBy == "name" {
			sort.Slice(items, func(i, j int) bool {
				if descending {
					return items[i].GetName() > items[j].GetName()
				} else {
					return items[i].GetName() < items[j].GetName()
				}
			})
		} else if sortBy == "doc_num" {
			sort.Slice(items, func(i, j int) bool {
				if descending {
					return items[i].DocNum > items[j].DocNum
				} else {
					return items[i].DocNum < items[j].DocNum
				}
			})
		} else if sortBy == "shard_num" {
			sort.Slice(items, func(i, j int) bool {
				if descending {
					return items[i].ShardNum > items[j].ShardNum
				} else {
					return items[i].ShardNum < items[j].ShardNum
				}
			})
		} else if sortBy == "storage_size" {
			sort.Slice(items, func(i, j int) bool {
				if descending {
					return items[i].StorageSize > items[j].StorageSize
				} else {
					return items[i].StorageSize < items[j].StorageSize
				}
			})
		} else if sortBy == "storage_type" {
			sort.Slice(items, func(i, j int) bool {
				if descending {
					return items[i].StorageType > items[j].StorageType
				} else {
					return items[i].StorageType < items[j].StorageType
				}
			})
		}
	}

	page.Total = int64(len(items))
	startIndex, endIndex := page.GetStartEndIndex()
	if endIndex > 0 {
		items = items[startIndex:endIndex]
	} else {
		items = []*core.Index{}
	}

	c.JSON(http.StatusOK, gin.H{
		"list": items,
		"page": page,
	})
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
