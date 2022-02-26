package handlers

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/index"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc/pkg/core"
	"github.com/prabhatsharma/zinc/pkg/startup"
)

func BulkHandler(c *gin.Context) {
	target := c.Param("target")

	ret, err := BulkHandlerWorker(target, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "bulk data inserted", "record_count": ret.Count})
}

func ESBulkHandler(c *gin.Context) {
	target := c.Param("target")

	startTime := time.Now()
	ret, err := BulkHandlerWorker(target, c.Request.Body)
	ret.Took = int(time.Since(startTime) / time.Millisecond)
	if err != nil {
		ret.Error = err.Error()
	}
	c.JSON(http.StatusOK, ret)
}

func BulkHandlerWorker(target string, body io.ReadCloser) (*BulkResponse, error) {
	bulkRes := &BulkResponse{Items: []map[string]*BulkResponseItem{}}

	// Prepare to read the entire raw text of the body
	scanner := bufio.NewScanner(body)
	defer body.Close()

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

	for scanner.Scan() { // Read each line
		var doc map[string]interface{}
		err := json.Unmarshal(scanner.Bytes(), &doc) // Read each line as JSON and store it in doc
		if err != nil {
			log.Print(err)
		}

		// This will process the data line in the request. Each data line is preceded by a metadata line.
		// Docs at https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-bulk.html
		if nextLineIsData {
			nextLineIsData = false
			var docID = ""
			mintedID := false

			bulkRes.Count++

			if val, ok := lastLineMetaData["_id"]; ok && val != nil {
				docID = val.(string)
			}
			if docID == "" {
				docID = uuid.New().String()
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

			// Since this is a bulk request, we need to check if we already created a new batch for this index. We need to create 1 batch per index.
			if DoesExistInThisRequest(indexesInThisBatch, indexName) == -1 { // Add the list of indexes to the batch if it's not already there
				indexesInThisBatch = append(indexesInThisBatch, indexName)
				batch[indexName] = index.NewBatch()
			}

			_, exists := core.GetIndex(indexName)
			if !exists { // If the requested indexName does not exist then create it
				newIndex, err := core.NewIndex(indexName, "disk", core.UseNewIndexMeta)
				if err != nil {
					return bulkRes, err
				}
				// store index
				core.StoreIndex(newIndex)
			}

			bdoc, err := core.ZINC_INDEX_LIST[indexName].BuildBlugeDocumentFromJSON(docID, &doc)
			if err != nil {
				return bulkRes, err
			}

			documentsInBatch++

			// Add the documen to the batch. We will persist the batch to the index
			// when we have processed all documents in the request
			if !mintedID {
				batch[indexName].Update(bdoc.ID(), bdoc)
			} else {
				batch[indexName].Insert(bdoc)
			}

			if documentsInBatch >= batchSize {
				for _, indexN := range indexesInThisBatch {
					// Persist the batch to the index
					err := core.ZINC_INDEX_LIST[indexN].Writer.Batch(batch[indexN])
					if err != nil {
						log.Printf("Error updating batch: %v", err)
						return bulkRes, err
					}
					batch[indexN].Reset()
				}
				documentsInBatch = 0
			}

		} else { // This branch will process the metadata line in the request. Each metadata line is preceded by a data line.

			for k, v := range doc {
				if k == "index" || k == "create" || k == "update" {
					nextLineIsData = true

					lastLineMetaData["operation"] = k

					if _, ok := v.(map[string]interface{}); !ok {
						// return errors.New("bulk index data format error")
						continue
					}

					if v.(map[string]interface{})["_index"] != "" { // if index is specified in metadata then it overtakes the index in the query path
						lastLineMetaData["_index"] = v.(map[string]interface{})["_index"]
					} else {
						lastLineMetaData["_index"] = target
					}

					lastLineMetaData["_id"] = v.(map[string]interface{})["_id"]
				} else if k == "delete" {
					nextLineIsData = false

					lastLineMetaData["operation"] = k
					lastLineMetaData["_index"] = v.(map[string]interface{})["_index"]
					lastLineMetaData["_id"] = v.(map[string]interface{})["_id"]

					// delete
					indexName := lastLineMetaData["_index"].(string)
					bdoc := bluge.NewDocument(lastLineMetaData["_id"].(string))
					if DoesExistInThisRequest(indexesInThisBatch, indexName) == -1 {
						indexesInThisBatch = append(indexesInThisBatch, indexName)
						batch[indexName] = index.NewBatch()
					}
					batch[indexName].Delete(bdoc.ID())

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

	for _, indexN := range indexesInThisBatch {
		writer := core.ZINC_INDEX_LIST[indexN].Writer
		// Persist the batch to the index
		err := writer.Batch(batch[indexN])
		if err != nil {
			log.Printf("Error updating batch: %v", err)
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
