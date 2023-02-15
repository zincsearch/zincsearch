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

	"github.com/zinclabs/zincsearch/pkg/core"
	"github.com/zinclabs/zincsearch/pkg/meta"
)

// @Id Refresh
// @Summary Resfresh index
// @security BasicAuth
// @Tags    Index
// @Produce json
// @Param   index  path  string  true  "Index"
// @Success 200 {object} meta.HTTPResponse
// @Failure 400 {object} meta.HTTPResponseError
// @Router /api/index/{index}/refresh [post]
func Refresh(c *gin.Context) {
	indexName := c.Param("target")
	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: "index " + indexName + " does not exists"})
		return
	}
	if err := index.Reopen(); err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, meta.HTTPResponse{Message: "ok"})
}
