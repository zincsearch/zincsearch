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
	"bufio"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zincsearch/pkg/config"
	"github.com/zinclabs/zincsearch/pkg/core"
	"github.com/zinclabs/zincsearch/pkg/ider"
	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/pkg/zutils/json"
)

// Multi accept multiple line json documents
// support use field `_id` set document id
//
// @Id Multi
// @Summary Multi documents
// @security BasicAuth
// @Tags    Document
// @Accept  plain
// @Produce json
// @Param   index  path  string  true  "Index"
// @Param   query  body  string  true  "Query"
// @Success 200 {object} meta.HTTPResponseRecordCount
// @Failure 400 {object} meta.HTTPResponseError
// @Failure 500 {object} meta.HTTPResponseError
// @Router /api/{index}/_multi [post]
func Multi(c *gin.Context) {
	target := c.Param("target")
	if target == "" {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: "target is empty"})
		return
	}

	defer c.Request.Body.Close()
	count, err := MultiWorker(target, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, meta.HTTPResponseRecordCount{Message: "multiple data inserted", RecordCount: count})
}

func MultiWorker(indexName string, body io.Reader) (int64, error) {
	// Prepare to read the entire raw text of the body
	scanner := bufio.NewScanner(body)

	// Set 1 MB max per line. docs at - https://pkg.go.dev/bufio#pkg-constants
	// This is the max size of a line in a file that we will process
	maxCapacityPerLine := config.Global.MaxDocumentSize
	buf := make([]byte, maxCapacityPerLine)
	scanner.Buffer(buf, maxCapacityPerLine)

	var doc map[string]interface{}
	var err error
	var count int64
	newIndex, _, err := core.GetOrCreateIndex(indexName, "", 0)
	if err != nil {
		return count, err
	}

	for scanner.Scan() { // Read each line
		for k := range doc {
			delete(doc, k)
		}
		if err = json.Unmarshal(scanner.Bytes(), &doc); err != nil {
			log.Error().Msgf("multi.json.Unmarshal: %s, err %s", scanner.Text(), err.Error())
			continue
		}

		update := false

		var docID = ""
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

	if err := scanner.Err(); err != nil {
		return count, err
	}

	return count, nil
}
