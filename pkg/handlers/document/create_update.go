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
)

func CreateUpdate(c *gin.Context) {
	indexName := c.Param("target")
	docID := c.Param("id") // ID for the document to be updated provided in URL path

	var err error
	var doc map[string]interface{}
	if err = c.BindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mintedID := false

	// If id field is present then use it, else create a new UUID and use it
	if id, ok := doc["_id"]; ok {
		docID = id.(string)
	}
	if docID == "" {
		docID = ider.Generate()
		mintedID = true
	}

	// If the index does not exist, then create it
	index, exists := core.GetIndex(indexName)
	if !exists {
		// Create a new index with disk storage as default
		index, err = core.NewIndex(indexName, "disk", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// store index
		_ = core.StoreIndex(index)
	}

	err = index.UpdateDocument(docID, doc, mintedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": docID})
}
