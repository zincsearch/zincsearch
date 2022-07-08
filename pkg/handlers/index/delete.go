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
	"github.com/zinclabs/zinc/pkg/meta"
)

// Delete deletes a zinc index and its associated data. Be careful using thus as you ca't undo this action.

// @Id DeleteIndex
// @Summary Delete index
// @Tags    Index
// @Produce json
// @Param   index  path  string  true  "Index"
// @Success 200 {object} meta.HTTPResponseIndex
// @Failure 400 {object} meta.HTTPResponseError
// @Failure 500 {object} meta.HTTPResponseError
// @Router /api/index/{index} [delete]
func Delete(c *gin.Context) {
	indexNames := c.Param("target")
	if indexNames == "" {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: "index name cannot be empty"})
		return
	}

	for _, indexName := range strings.Split(indexNames, ",") {
		if err := core.DeleteIndex(indexName); err != nil {
			c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, meta.HTTPResponse{
		Message: "deleted",
	})
}
