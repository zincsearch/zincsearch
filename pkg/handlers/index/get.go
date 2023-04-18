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

	"github.com/gin-gonic/gin"

	"github.com/zincsearch/zincsearch/pkg/core"
	"github.com/zincsearch/zincsearch/pkg/meta"
)

// @Id GetIndex
// @Summary Get index metadata
// @security BasicAuth
// @Tags    Index
// @Produce json
// @Param   index  path  string  true  "Index"
// @Success 200 {object} core.Index
// @Failure 404 {object} meta.HTTPResponseError
// @Router /api/index/{index} [get]
func Get(c *gin.Context) {
	indexName := c.Param("target")
	if indexName == "" {
		c.JSON(http.StatusNotFound, meta.HTTPResponseError{Error: "index " + indexName + " does not exists"})
		return
	}

	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusNotFound, meta.HTTPResponseError{Error: "index " + indexName + " does not exists"})
		return
	}

	c.JSON(http.StatusOK, index)
}
