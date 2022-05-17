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
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/index"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/ider"
	"github.com/zinclabs/zinc/pkg/startup"
)

func Bulk(c *gin.Context) {
	target := c.Param("target")

	defer c.Request.Body.Close()

	ret, err := BulkWorker(target, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "bulk data inserted", "record_count": ret.Count})
}

func ESBulk(c *gin.Context) {
	target := c.Param("target")

	defer c.Request.Body.Close()

	startTime := time.Now()
	ret, err := BulkWorker(target, c.Request.Body)
	ret.Took = int(time.Since(startTime) / time.Millisecond)
	if err != nil {
		ret.Error = err.Error()
	}
	c.JSON(http.StatusOK, ret)
}

func BulkWorker(target string, body io.Reader) (*BulkResponse, error) {
	bulkRes := &BulkResponse{Items: []map[string]*BulkResponseItem{}}

	// Prepare to read the entire raw text of the body
	scanner := bufio.NewScanner(body)

	// force set batchSize
	batchSize := startup.LoadBatchSize()

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
			log.Error().Msgf("bulk.json.Unmarshal: err %s", err.Error())
			continue
		}

		// This will process the data line in the request. Each data line is preceded by a metadata line.
		// Docs at https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-bulk.html
		if nextLineIsData {
			bulkRes.Count++
			nextLineIsData = false
			mintedID := false

			var docID = ""
			if val, ok := lastLineMetaData["_id"]; ok && val != nil {
				docID = val.(string)
			}
			if docID == "" {
				docID = ider.Generate()
				mintedID = true
			}

			indexName := lastLineMetaData["_index"].(string)
			operation := lastLineMetaData["operation"].(string)
			switch operation {
			case "index":
				bulkRes.Items = append(bulkRes.Items, map[string]*BulkResponseItem{
					"index": NewBulkResponseItem(bulkRes.Count, indexName, docID, "created", nil),
				})
			case "create":
				bulkRes.Items = append(bulkRes.Items, map[string]*BulkResponseItem{
					"index": NewBulkResponseItem(bulkRes.Count, indexName, docID, "created", nil),
				})
			case "update":
				bulkRes.Items = append(bulkRes.Items, map[string]*BulkResponseItem{
					"index": NewBulkResponseItem(bulkRes.Count, indexName, docID, "updated", nil),
				})
			default:
			}

			_, exists := core.GetIndex(indexName)
			if !exists { // If the requested indexName does not exist then create it
				newIndex, err := core.NewIndex(indexName, "disk", nil)
				if err != nil {
					return bulkRes, err
				}
				// store index
				if err := core.StoreIndex(newIndex); err != nil {
					return bulkRes, err
				}
			}

			// Since this is a bulk request, we need to check if we already created a new batch for this index. We need to create 1 batch per index.
			if DoesExistInThisRequest(indexesInThisBatch, indexName) == -1 { // Add the list of indexes to the batch if it's not already there
				indexesInThisBatch = append(indexesInThisBatch, indexName)
				batch[indexName] = index.NewBatch()
			}

			bdoc, err := core.ZINC_INDEX_LIST[indexName].BuildBlugeDocumentFromJSON(docID, doc)
			if err != nil {
				return bulkRes, err
			}

			// Add the documen to the batch. We will persist the batch to the index
			// when we have processed all documents in the request
			if !mintedID {
				batch[indexName].Update(bdoc.ID(), bdoc)
			} else {
				batch[indexName].Insert(bdoc)
			}

			documentsInBatch++

			// refresh index stats
			core.ZINC_INDEX_LIST[indexName].GainDocsCount(1)

			if documentsInBatch >= batchSize {
				for _, indexName := range indexesInThisBatch {
					// Persist the batch to the index
					if err := core.ZINC_INDEX_LIST[indexName].Writer.Batch(batch[indexName]); err != nil {
						log.Error().Msgf("bulk: index updating batch err %s", err.Error())
						return bulkRes, err
					}
					batch[indexName].Reset()
				}
				documentsInBatch = 0
			}

		} else { // This branch will process the metadata line in the request. Each metadata line is preceded by a data line.

			for k, v := range doc {
				if k == "index" || k == "create" || k == "update" {
					nextLineIsData = true

					lastLineMetaData["operation"] = k

					if _, ok := v.(map[string]interface{}); !ok {
						return nil, errors.New("bulk index data format error")
					}

					if v.(map[string]interface{})["_index"] != "" { // if index is specified in metadata then it overtakes the index in the query path
						lastLineMetaData["_index"] = v.(map[string]interface{})["_index"]
					} else {
						lastLineMetaData["_index"] = target
					}
					if lastLineMetaData["_index"] == "" {
						return nil, errors.New("bulk index data format error")
					}

					lastLineMetaData["_id"] = v.(map[string]interface{})["_id"]
				} else if k == "delete" {
					nextLineIsData = false

					lastLineMetaData["operation"] = k
					lastLineMetaData["_id"] = v.(map[string]interface{})["_id"]
					if v.(map[string]interface{})["_index"] != "" { // if index is specified in metadata then it overtakes the index in the query path
						lastLineMetaData["_index"] = v.(map[string]interface{})["_index"]
					} else {
						lastLineMetaData["_index"] = target
					}
					if lastLineMetaData["_index"] == "" {
						return nil, errors.New("bulk index data format error")
					}

					// delete
					indexName := lastLineMetaData["_index"].(string)
					bdoc := bluge.NewDocument(lastLineMetaData["_id"].(string))
					if DoesExistInThisRequest(indexesInThisBatch, indexName) == -1 {
						indexesInThisBatch = append(indexesInThisBatch, indexName)
						batch[indexName] = index.NewBatch()
					}
					batch[indexName].Delete(bdoc.ID())
					core.ZINC_INDEX_LIST[indexName].ReduceDocsCount(1)

					bulkRes.Count++
					bulkRes.Items = append(bulkRes.Items, map[string]*BulkResponseItem{
						"delete": NewBulkResponseItem(bulkRes.Count, indexName, bdoc.ID().Field(), "deleted", nil),
					})
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return bulkRes, err
	}

	for _, indexName := range indexesInThisBatch {
		// Persist the batch to the index
		if err := core.ZINC_INDEX_LIST[indexName].Writer.Batch(batch[indexName]); err != nil {
			log.Printf("bulk: index updating batch err %s", err.Error())
			return bulkRes, err
		}
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

func NewBulkResponseItem(seqNo int, index, id, result string, err error) *BulkResponseItem {
	return &BulkResponseItem{
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
		Status: 200,
		SeqNo:  seqNo,
		Error:  err,
	}
}

type BulkResponse struct {
	Took   int                            `json:"took"`
	Errors bool                           `json:"errors"`
	Error  string                         `json:"error,omitempty"`
	Count  int                            `json:"count"`
	Items  []map[string]*BulkResponseItem `json:"items"`
}

type BulkResponseItem struct {
	Index   string                `json:"_index"`
	Type    string                `json:"_type"`
	ID      string                `json:"_id"`
	Version int64                 `json:"_version"`
	Result  string                `json:"result"`
	Shards  BulkResponseItemShard `json:"_shards"`
	Status  int                   `json:"status"`
	SeqNo   int                   `json:"seq_no"`
	Error   error                 `json:"error,omitempty"`
}

type BulkResponseItemShard struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}
