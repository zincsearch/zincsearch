/* Copyright 2022 Zinc Labs Inc and Contributors
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
	"github.com/zinclabs/zincsearch/pkg/ider"
	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/pkg/zutils"
)

// Bulkv2 accept JSONIngest json documents. Its a simpler and standard format to ingest data.
// support use field `_id` set document id
//
// @Id Bulkv2
// @Summary Bulkv2 documents
// @security BasicAuth
// @Tags    Document
// @Accept  json
// @Produce json
// @Param   query  body  meta.JSONIngest  true  "Query"
// @Success 200 {object} meta.HTTPResponseRecordCount
// @Failure 400 {object} meta.HTTPResponseError
// @Failure 500 {object} meta.HTTPResponseError
// @Router /api/_bulkv2 [post]
func Bulkv2(c *gin.Context) {
	target := c.Param("target")

	var body meta.JSONIngest

	if err := zutils.GinBindJSON(c, &body); err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	if target == "" {
		target = body.Index
	}

	defer c.Request.Body.Close()
	count, err := Bulkv2Worker(target, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, meta.HTTPResponseRecordCount{Message: "v2 data inserted", RecordCount: count})
}

// Bulkv2Worker accept JSONIngest json documents. It provides a simpler format to ingest data.
func Bulkv2Worker(indexName string, body meta.JSONIngest) (int64, error) {
	var err error
	var count int64
	newIndex, _, err := core.GetOrCreateIndex(indexName, "", 0)
	if err != nil {
		return count, err
	}

	for _, doc := range body.Records { // Read each line
		update := false

		docID := ""
		if val, ok := doc["_id"]; ok && val != nil {
			docID = val.(string)
		}
		if docID == "" {
			docID = ider.Generate()
		} else {
			update = true
		}

		err = newIndex.CreateDocument(docID, doc, update)
		if err != nil {
			return count, err
		}

		count++
	}

	return count, nil
}
