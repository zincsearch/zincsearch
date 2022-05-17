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
)

func Delete(c *gin.Context) {
	indexName := c.Param("target")
	docID := c.Param("id")

	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index does not exists"})
		return
	}

	bdoc := bluge.NewDocument(docID)
	err := index.Writer.Delete(bdoc.ID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	core.ZINC_INDEX_LIST[indexName].ReduceDocsCount(1)
	c.JSON(http.StatusOK, gin.H{"message": "deleted", "index": indexName, "id": docID})
}
