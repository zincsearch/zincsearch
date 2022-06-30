package benchmark

import (
	"bufio"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/zinclabs/zinc/pkg/core"
)

func prepareOlympicsData() (*core.Index, []map[string]interface{}) {
	f := getOlympicsFile()

	documents := make([]map[string]interface{}, 0, 100)
	scanner := bufio.NewScanner(f)

	scanner.Scan()
	scanner.Text()

	for scanner.Scan() {
		line := scanner.Text()
		var doc map[string]interface{}
		json.Unmarshal([]byte(line), &doc)
		documents = append(documents, doc)
	}

	index, has := core.GetIndex("olympics")
	if !has {
		var err error
		index, err = core.NewIndex("olympics", "disk", nil)
		if err != nil {
			panic(err)
		}
		err = core.StoreIndex(index)
		if err != nil {
			panic(err)
		}
	}

	return index, documents
}

type CRUDFunctions struct {
	create func(string, map[string]interface{}) error
	update func(string, map[string]interface{}) error
	delete func(string) error
}

const documentNumber = 300

func runCRUDOperations(b *testing.B, documents []map[string]interface{}, fns CRUDFunctions) {
	documents = documents[0:documentNumber]
	b.StartTimer()

	for i := range documents {
		fns.update(strconv.Itoa(i), documents[i])
	}

	for i := range documents {
		if i%2 == 0 {
			continue
		}
		fns.update(strconv.Itoa(i), documents[i])
	}

	for i := range documents {
		if i%3 != 0 {
			continue
		}
		fns.delete(strconv.Itoa(i))
	}
}

func BenchmarkWal(b *testing.B) {
	b.SetParallelism(1)
	index, data := prepareOlympicsData()
	runCRUDOperations(b, data, CRUDFunctions{
		create: index.CreateDocumentAsync,
		update: index.UpdateDocumentAsync,
		delete: index.DeleteDocumentAsync,
	})

	expectedDocs := documentNumber - documentNumber/3
	for {
		index.UpdateMetadata()
		if index.DocNum == uint64(expectedDocs) {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func BenchmarkNoWal(b *testing.B) {
	index, data := prepareOlympicsData()
	runCRUDOperations(b, data, CRUDFunctions{
		create: func(id string, data map[string]interface{}) error {
			return index.CreateDocument(id, data, false)
		},
		update: index.UpdateDocument,
		delete: index.DeleteDocument,
	})
}
