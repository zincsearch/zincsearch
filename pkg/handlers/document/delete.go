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

	"github.com/blugelabs/bluge"
	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/meta"
)

// @Id DeleteDocument
// @Summary Delete document
// @Tags    Document
// @Param   index  path  string  true  "Index"
// @Param   id     path  string  true  "ID"
// @Success 200 {object} meta.HTTPResponseDocument
// @Failure 400 {object} meta.HTTPResponseError
// @Failure 500 {object} meta.HTTPResponseError
// @Router /api/{index}/_doc/{id} [delete]
func Delete(c *gin.Context) {
	indexName := c.Param("target")
	docID := c.Param("id")

	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: "index does not exists"})
		return
	}

	bdoc := bluge.NewDocument(docID)
	writers, err := index.GetWriters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponseError{Error: err.Error()})
		return
	}
	for _, w := range writers {
		err = w.Delete(bdoc.ID())
		if err != nil {
			c.JSON(http.StatusInternalServerError, meta.HTTPResponseError{Error: err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, meta.HTTPResponseDocument{Message: "deleted", Index: indexName, ID: docID})
}
