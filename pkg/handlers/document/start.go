package document

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-json"
)

func init() {
	fmt.Println("document package init")

	go IngestFromWAL()
}

// IngestFromWAL reads from the WAL and ingests the documents into the index
func IngestFromWAL() {
	SEPARATOR := []byte{0xff}
	var messageType, docID, indexName string
	var message []byte

	fmt.Println("document package ingestFromWAL")
	for {
		wals := ReadFromWAL(1) // Read 1 record from WAL. ReadFromWAL returns a [][]byte
		if len(wals) == 0 {
			continue
		}
		entry := wals[0]
		fmt.Println("entry: ", entry)

		// ingest the entry

		record := bytes.Split(entry, SEPARATOR)
		messageType = string(record[0])
		docID = string(record[1])
		indexName = string(record[2])
		message = record[3]

		if messageType == "single" {
			// ingest the message
			fmt.Println("messageType: ", messageType)
			var doc map[string]interface{}
			err := json.Unmarshal(message, &doc)
			if err != nil {
				fmt.Println("err: ", err)
			}
			createUpdateDocumentWorker(doc, docID, indexName) // Ingest the message
			MarkWALAsConsumed(1)                              // Mark the message as consumed
		} else if messageType == "bulk" {
			// ingest the message
			fmt.Println("messageType: ", messageType)
		}
	}
}
