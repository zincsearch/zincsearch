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

package document

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zincsearch/pkg/core"
	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/pkg/zutils"
)

// @Id GetDocument
// @Summary get document with id
// @security BasicAuth
// @Tags    Document
// @Accept  json
// @Produce json
// @Param   index  path  string  true  "Index"
// @Param   id     path  string  true  "ID"
// @Success 200 {object} meta.Hit
// @Failure 400 {object} meta.HTTPResponseError
// @Failure 500 {object} meta.HTTPResponseError
// @Router /api/{index}/_doc/{id} [get]
func Get(c *gin.Context) {
	indexName := c.Param("target")
	docID := c.Param("id")

	if docID == "" {
		zutils.GinRenderJSON(c, http.StatusBadRequest, meta.HTTPResponseError{Error: "id is empty"})
		return
	}

	// If the index does not exist, then create it
	index, _, err := core.GetOrCreateIndex(indexName, "", 0)
	if err != nil {
		zutils.GinRenderJSON(c, http.StatusInternalServerError, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	source, err := index.GetDocument(docID)
	if err != nil {
		zutils.GinRenderJSON(c, http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}
	zutils.GinRenderJSON(c, http.StatusOK, source)
}
