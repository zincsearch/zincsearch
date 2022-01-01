package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"

	"github.com/blugelabs/bluge/index"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prabhatsharma/zinc/pkg/core"
)

func BulkHandler(c *gin.Context) {
	target := c.Param("target")

	body := c.Request.Body

	err := BulkHandlerWorker(target, &body)

	if err != nil {
		c.JSON(200, gin.H{
			"message": err,
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "bulk data inserted",
	})
}

func BulkHandlerWorker(target string, body *io.ReadCloser) error {
	// Prepare to read the entire raw text of the body
	scanner := bufio.NewScanner(*body)

	// Set 1 MB max per line. docs at - https://pkg.go.dev/bufio#pkg-constants
	// This is the max size of a line in a file that we will process
	const maxCapacityPerLine = 1024 * 1024
	buf := make([]byte, maxCapacityPerLine)
	scanner.Buffer(buf, maxCapacityPerLine)

	nextLineIsData := false
	lastLineMetaData := make(map[string]interface{})

	batch := make(map[string]*index.Batch)
	var indexesInThisBatch []string

	documentsInBatch := 0

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
			var id = ""

			if val, ok := lastLineMetaData["id"]; ok {
				id = val.(string)
			} else {
				id = uuid.New().String()
			}

			indexName := lastLineMetaData["_index"].(string)

			// Since this is a bulk request, we need to check if we already created a new batch for this index. We need to create 1 batch per index.
			if DoesExistInThisRequest(indexesInThisBatch, indexName) == -1 { // Add the list of indexes to the batch if it's not already there
				indexesInThisBatch = append(indexesInThisBatch, indexName)
				batch[indexName] = index.NewBatch()
			}

			exists, _ := core.IndexExists(indexName)

			if !exists { // If the requested indexName does not exist then create it
				newIndex, err := core.NewIndex(indexName, "disk")

				if err != nil {
					return err
				}

				core.ZINC_INDEX_LIST[indexName] = newIndex // Load the index in memory
			}

			bdoc, err := core.ZINC_INDEX_LIST[indexName].BuildBlugeDocumentFromJSON(id, &doc)

			if err != nil {
				return err
			}

			documentsInBatch++

			// Add the documen to the batch. We will persist the batch to the index
			// when we have processed all documents in the request
			batch[indexName].Update(bdoc.ID(), bdoc)

		} else { // This branch will process the metadata line in the request. Each metadata line is preceded by a data line.

			for k, v := range doc {
				if k == "index" || k == "create" || k == "update" {
					nextLineIsData = true

					lastLineMetaData["operation"] = k

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

				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for _, indexN := range indexesInThisBatch {

		writer := core.ZINC_INDEX_LIST[indexN].Writer

		// Persist the batch to the index
		err := writer.Batch(batch[indexN])
		if err != nil {
			fmt.Println("Error updating batch: ", err.Error())
			log.Print(err)
			return err
		}
	}

	return nil

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
