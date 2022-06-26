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

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/index"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/ider"
	"github.com/zinclabs/zinc/pkg/meta"
)

// @Id Bulk
// @Summary Bulk documents
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

	indexes, ret, err := BulkWorker(target, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	// check shards
	for name := range indexes {
		index, ok := core.ZINC_INDEX_LIST.Get(name)
		if !ok {
			continue
		}
		_ = index.CheckShards()
	}

	c.JSON(http.StatusOK, meta.HTTPResponseRecordCount{Message: "bulk data inserted", RecordCount: ret.Count})
}

// @Id ESBulk
// @Summary ES bulk documents
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

	startTime := time.Now()
	indexes, ret, err := BulkWorker(target, c.Request.Body)
	ret.Took = int(time.Since(startTime) / time.Millisecond)
	if err != nil {
		ret.Error = err.Error()
	}

	// check shards
	for name := range indexes {
		index, ok := core.ZINC_INDEX_LIST.Get(name)
		if !ok {
			continue
		}
		_ = index.CheckShards()
	}

	// update seqNo
	atomic.AddInt64(&globalSeqNo, int64(ret.Count))

	c.JSON(http.StatusOK, ret)
}

func BulkWorker(target string, body io.Reader) (map[string]struct{}, *BulkResponse, error) {
	bulkRes := &BulkResponse{Items: []map[string]BulkResponseItem{}}
	bulkIndexes := make(map[string]struct{})

	// Prepare to read the entire raw text of the body
	scanner := bufio.NewScanner(body)

	// force set batchSize
	batchSize := config.Global.BatchSize

	// Set 1 MB max per line. docs at - https://pkg.go.dev/bufio#pkg-constants
	// This is the max size of a line in a file that we will process
	const maxCapacityPerLine = 1024 * 1024
	buf := make([]byte, maxCapacityPerLine)
	scanner.Buffer(buf, maxCapacityPerLine)

	nextLineIsData := false
	lastLineMetaData := make(map[string]interface{})

	batch := make(map[string]*index.Batch)
	var indexesInThisBatch []string
	var documentsInBatch int
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

			var docID = ""
			if val, ok := lastLineMetaData["_id"]; ok && val != nil {
				docID = val.(string)
			}
			if docID == "" {
				docID = ider.Generate()
			} else {
				update = true
			}

			indexName := lastLineMetaData["_index"].(string)
			operation := lastLineMetaData["operation"].(string)
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

			_, exists := core.GetIndex(indexName)
			if !exists { // If the requested indexName does not exist then create it
				newIndex, err := core.NewIndex(indexName, "", nil)
				if err != nil {
					return bulkIndexes, bulkRes, err
				}
				// store index
				if err := core.StoreIndex(newIndex); err != nil {
					return bulkIndexes, bulkRes, err
				}
			}
			bulkIndexes[indexName] = struct{}{}

			// Since this is a bulk request, we need to check if we already created a new batch for this index. We need to create 1 batch per index.
			if DoesExistInThisRequest(indexesInThisBatch, indexName) == -1 { // Add the list of indexes to the batch if it's not already there
				indexesInThisBatch = append(indexesInThisBatch, indexName)
				batch[indexName] = index.NewBatch()
			}

			newIndex, _ := core.GetIndex(indexName)
			bdoc, err := newIndex.BuildBlugeDocumentFromJSON(docID, doc)
			if err != nil {
				return bulkIndexes, bulkRes, err
			}

			// Add the documen to the batch. We will persist the batch to the index
			// when we have processed all documents in the request
			if update {
				batch[indexName].Update(bdoc.ID(), bdoc)
			} else {
				batch[indexName].Insert(bdoc)
			}

			documentsInBatch++

			if documentsInBatch >= batchSize {
				for _, indexName := range indexesInThisBatch {
					// Persist the batch to the index
					newIndex, _ := core.GetIndex(indexName)
					writer, err := newIndex.GetWriter()
					if err != nil {
						log.Error().Msgf("bulk: index updating batch err %s", err.Error())
						return bulkIndexes, bulkRes, err
					}
					if err := writer.Batch(batch[indexName]); err != nil {
						log.Error().Msgf("bulk: index updating batch err %s", err.Error())
						return bulkIndexes, bulkRes, err
					}
					batch[indexName].Reset()
				}
				documentsInBatch = 0
			}

		} else { // This branch will process the metadata line in the request. Each metadata line is preceded by a data line.

			for k, v := range doc {
				vm, ok := v.(map[string]interface{})
				if !ok {
					return nil, nil, errors.New("bulk index data format error")
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
						return nil, nil, errors.New("bulk index data format error")
					}
					lastLineMetaData["_id"] = vm["_id"]
				} else if k == "delete" {
					nextLineIsData = false

					lastLineMetaData["operation"] = k
					lastLineMetaData["_id"] = vm["_id"]
					if vm["_index"] != "" { // if index is specified in metadata then it overtakes the index in the query path
						lastLineMetaData["_index"] = vm["_index"]
					} else {
						lastLineMetaData["_index"] = target
					}
					if lastLineMetaData["_index"] == "" {
						return nil, nil, errors.New("bulk index data format error")
					}

					// delete
					indexName := lastLineMetaData["_index"].(string)
					bdoc := bluge.NewDocument(lastLineMetaData["_id"].(string))
					if DoesExistInThisRequest(indexesInThisBatch, indexName) == -1 {
						indexesInThisBatch = append(indexesInThisBatch, indexName)
						batch[indexName] = index.NewBatch()
					}
					batch[indexName].Delete(bdoc.ID())

					bulkRes.Count++
					bulkRes.Items = append(bulkRes.Items, map[string]BulkResponseItem{
						"delete": NewBulkResponseItem(bulkRes.Count, indexName, bdoc.ID().Field(), "deleted", nil),
					})
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return bulkIndexes, bulkRes, err
	}

	for _, indexName := range indexesInThisBatch {
		// Persist the batch to the index
		newIndex, _ := core.GetIndex(indexName)
		writer, err := newIndex.GetWriter()
		if err != nil {
			log.Error().Msgf("bulk: index updating batch err %s", err.Error())
			return bulkIndexes, bulkRes, err
		}
		if err := writer.Batch(batch[indexName]); err != nil {
			log.Printf("bulk: index updating batch err %s", err.Error())
			return bulkIndexes, bulkRes, err
		}
	}

	return bulkIndexes, bulkRes, nil
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
