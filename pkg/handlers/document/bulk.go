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
	"errors"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/zincsearch/zincsearch/pkg/config"
	"github.com/zincsearch/zincsearch/pkg/core"
	"github.com/zincsearch/zincsearch/pkg/ider"
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/zutils"
	"github.com/zincsearch/zincsearch/pkg/zutils/json"
)

// Bulk accept multiple documents, first line index metadata, second line document
//
// @Id Bulk
// @Summary Bulk documents
// @security BasicAuth
// @Tags    Document
// @Accept  plain
// @Produce json
// @Param   query  body  string  true  "Query"
// @Success 200 {object} meta.HTTPResponseRecordCount
// @Failure 500 {object} meta.HTTPResponseError
// @Router /api/_bulk [post]
func Bulk(c *gin.Context) {
	target := c.Param("target")

	defer c.Request.Body.Close()

	ret, err := BulkWorker(target, c.Request.Body)
	if err != nil {
		zutils.GinRenderJSON(c, http.StatusInternalServerError, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	zutils.GinRenderJSON(c, http.StatusOK, meta.HTTPResponseRecordCount{Message: "bulk data inserted", RecordCount: ret.Count})
}

// ESBulk accept multiple documents, first line index metadata, second line document
//
// @Id ESBulk
// @Summary ES bulk documents
// @security BasicAuth
// @Tags    Document
// @Accept  plain
// @Produce json
// @Param   query  body  string  true  "Query"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} meta.HTTPResponseError
// @Router /es/_bulk [post]
func ESBulk(c *gin.Context) {
	target := c.Param("target")

	defer c.Request.Body.Close()

	ret, err := BulkWorker(target, c.Request.Body)
	if err != nil {
		ret.Error = err.Error()
	}

	startTime := time.Now()
	ret.Took = int(time.Since(startTime) / time.Millisecond)
	// update seqNo
	atomic.AddInt64(&globalSeqNo, int64(ret.Count))

	zutils.GinRenderJSON(c, http.StatusOK, ret)
}

func BulkWorker(target string, body io.Reader) (*BulkResponse, error) {
	bulkRes := &BulkResponse{Items: []map[string]BulkResponseItem{}}

	// Prepare to read the entire raw text of the body
	scanner := bufio.NewScanner(body)

	// Set 1 MB max per line. docs at - https://pkg.go.dev/bufio#pkg-constants
	// This is the max size of a line in a file that we will process
	maxCapacityPerLine := config.Global.MaxDocumentSize
	buf := make([]byte, maxCapacityPerLine)
	scanner.Buffer(buf, maxCapacityPerLine)

	nextLineIsData := false
	lastLineMetaData := make(map[string]interface{})

	var doc map[string]interface{}
	var err error
	for scanner.Scan() { // Read each line
		for k := range doc {
			delete(doc, k)
		}
		if err = json.Unmarshal(scanner.Bytes(), &doc); err != nil {
			log.Error().Msgf("bulk.json.Unmarshal: %s, err %s", scanner.Text(), err.Error())
			continue
		}

		// This will process the data line in the request. Each data line is preceded by a metadata line.
		// Docs at https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-bulk.html
		if nextLineIsData {
			bulkRes.Count++
			nextLineIsData = false
			update := false

			docID := ""
			if val, ok := lastLineMetaData["_id"]; ok && val != nil {
				docID = val.(string)
			}
			if docID == "" {
				docID = ider.Generate()
			} else {
				update = true
			}
			suppliedIndexName := lastLineMetaData["_index"]
			if suppliedIndexName == nil && len(target) > 0 {
				suppliedIndexName = target
			}
			suppliedOperation := lastLineMetaData["operation"]
			if suppliedOperation == nil && len(target) > 0 {
				suppliedOperation = "create"
			}
			indexName := suppliedIndexName.(string)
			operation := suppliedOperation.(string)
			switch operation {
			case "index":
				bulkRes.Items = append(bulkRes.Items, map[string]BulkResponseItem{
					"index": NewBulkResponseItem(bulkRes.Count, indexName, docID, "created", nil),
				})
			case "create":
				bulkRes.Items = append(bulkRes.Items, map[string]BulkResponseItem{
					"index": NewBulkResponseItem(bulkRes.Count, indexName, docID, "created", nil),
				})
			case "update":
				bulkRes.Items = append(bulkRes.Items, map[string]BulkResponseItem{
					"index": NewBulkResponseItem(bulkRes.Count, indexName, docID, "updated", nil),
				})
			default:
			}

			newIndex, _, err := core.GetOrCreateIndex(indexName, "", 0)
			if err != nil {
				return bulkRes, err
			}

			err = newIndex.CreateDocument(docID, doc, update)
			if err != nil {
				return bulkRes, err
			}

		} else { // This branch will process the metadata line in the request. Each metadata line is preceded by a data line.

			for k, v := range doc {
				vm, ok := v.(map[string]interface{})
				if !ok {
					return nil, errors.New("bulk index data format error")
				}
				for k := range lastLineMetaData {
					delete(lastLineMetaData, k)
				}
				if k == "index" || k == "create" || k == "update" {
					nextLineIsData = true
					lastLineMetaData["operation"] = k

					if vm["_index"] != "" { // if index is specified in metadata then it overtakes the index in the query path
						lastLineMetaData["_index"] = vm["_index"]
					} else {
						lastLineMetaData["_index"] = target
					}
					if lastLineMetaData["_index"] == "" {
						return nil, errors.New("bulk index data format error")
					}
					lastLineMetaData["_id"] = vm["_id"]
				} else if k == "delete" {
					nextLineIsData = false
					docID := vm["_id"].(string)
					indexName := target
					if vm["_index"] != "" { // if index is specified in metadata then it overtakes the index in the query path
						indexName = vm["_index"].(string)
					}
					if indexName == "" {
						return nil, errors.New("bulk index data format error")
					}

					newIndex, _, err := core.GetOrCreateIndex(indexName, "", 0)
					if err != nil {
						return bulkRes, err
					}

					// delete
					err = newIndex.DeleteDocument(docID)
					bulkRes.Count++
					bulkRes.Items = append(bulkRes.Items, map[string]BulkResponseItem{
						"delete": NewBulkResponseItem(bulkRes.Count, indexName, docID, "deleted", err),
					})
				} else {
					lastLineMetaData["_index"] = target
					lastLineMetaData["operation"] = "index"
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return bulkRes, err
	}

	return bulkRes, nil
}

// DoesExistInThisRequest takes a slice and looks for an element in it. If found it will
// return it's index, otherwise it will return -1.
func DoesExistInThisRequest(slice []string, val string) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}

func NewBulkResponseItem(seqNo int64, index, id, result string, err error) BulkResponseItem {
	return BulkResponseItem{
		Index:   index,
		Type:    "_doc",
		ID:      id,
		Version: 1,
		Result:  result,
		Shards: BulkResponseItemShard{
			Total:      1,
			Successful: 1,
			Failed:     0,
		},
		Status:      200,
		SeqNo:       globalSeqNo + seqNo,
		PrimaryTerm: 1,
		Error:       err,
	}
}

var globalSeqNo int64

type BulkResponse struct {
	Took   int                           `json:"took"`
	Errors bool                          `json:"errors"`
	Error  string                        `json:"error,omitempty"`
	Items  []map[string]BulkResponseItem `json:"items"`
	Count  int64                         `json:"-"`
}

type BulkResponseItem struct {
	Index       string                `json:"_index"`
	Type        string                `json:"_type"`
	ID          string                `json:"_id"`
	Version     int64                 `json:"_version"`
	Result      string                `json:"result"`
	Status      int                   `json:"status"`
	Shards      BulkResponseItemShard `json:"_shards"`
	SeqNo       int64                 `json:"_seq_no"`
	PrimaryTerm int                   `json:"_primary_term"`
	Error       error                 `json:"error,omitempty"`
}

type BulkResponseItemShard struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}
