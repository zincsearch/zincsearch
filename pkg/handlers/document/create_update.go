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
	"github.com/zinclabs/zinc/pkg/ider"
	"github.com/zinclabs/zinc/pkg/meta"
)

// @Summary Create update document
// @Tags    Document
// @Param   target    path  string  true  "Index"
// @Param   id        path  string  false "ID"
// @Param   document  body  map[string]interface{}  true  "Document"
// @Success 200 {object} meta.HTTPResponse
// @Failure 400 {object} meta.HTTPResponse
// @Failure 500 {object} meta.HTTPResponse
// @Router /api/:target/_doc [post]
// @Router /api/:target/_doc/:id [put]
func CreateUpdate(c *gin.Context) {
	indexName := c.Param("target")
	docID := c.Param("id") // ID for the document to be updated provided in URL path

	var err error
	var doc map[string]interface{}
	if err = c.BindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: err.Error()})
		return
	}

	err = sendToWAL("single", docID, indexName, &c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponse{Error: err.Error()})
		return
	}

	// err = createUpdateDocumentWorker(doc, docID, indexName)

	if err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, meta.HTTPResponse{Message: "ok", ID: docID})
}

// createUpdateDocument is a helper function to create or update a document
func createUpdateDocumentWorker(doc map[string]interface{}, docID string, indexName string) error {

	update := false
	// If id field is present then use it, else create a new UUID and use it
	if id, ok := doc["_id"]; ok {
		docID = id.(string)
	}
	if docID == "" {
		docID = ider.Generate()
	} else {
		update = true
	}

	var err error

	// If the index does not exist, then create it
	index, exists := core.GetIndex(indexName)
	if !exists {
		// Create a new index with disk storage as default
		index, err = core.NewIndex(indexName, "disk", nil)
		if err != nil {
			return err
			// c.JSON(http.StatusInternalServerError, meta.HTTPResponse{Error: err.Error()})
			// return
		}
		// store index
		_ = core.StoreIndex(index)
	}

	err = index.CreateDocument(docID, doc, update)
	if err != nil {
		return err
	}

	// check shards
	_ = index.CheckShards()

	return nil
}
