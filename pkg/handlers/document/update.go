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

	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/meta"
)

// @Summary Update document with id
// @Tags    Document
// @Param   target path  string  true  "Index"
// @Param   id     path  string  true  "ID"
// @Param   document  body  map[string]interface{}  true  "Document"
// @Success 200 {object} meta.HTTPResponse
// @Failure 400 {object} meta.HTTPResponse
// @Failure 500 {object} meta.HTTPResponse
// @Router /api/:target/_update/:id [post]
func Update(c *gin.Context) {
	indexName := c.Param("target")
	docID := c.Param("id") // ID for the document to be updated provided in URL path

	var err error
	var doc map[string]interface{}
	if err = c.BindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: err.Error()})
		return
	}

	// If id field is present then use it, else create a new UUID and use it
	if id, ok := doc["_id"]; ok {
		docID = id.(string)
	}
	if docID == "" {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: "_id field is required"})
		return
	}

	// If the index does not exist, then create it
	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: "index does not exists"})
		return
	}

	err = index.UpdateDocument(docID, doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, meta.HTTPResponse{Message: "ok", ID: docID})
}
