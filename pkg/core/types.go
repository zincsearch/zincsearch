package core

import (
	"github.com/blugelabs/bluge"
)

// Nothing to handle in the error. If you can't load indexes then everything is broken.
var ZINC_INDEX_LIST map[string]*Index

var ZINC_SYSTEM_INDEX_LIST, _ = LoadZincSystemIndexes()

func init() {
	ZINC_INDEX_LIST, _ = LoadZincIndexesFromDisk()
	s3List, _ := LoadZincIndexesFromS3()
	// Load the indexes from disk.
	for k, v := range s3List {
		ZINC_INDEX_LIST[k] = v
	}
}

type Index struct {
	Name          string            `json:"name"`
	Writer        *bluge.Writer     `json:"blugeWriter"`
	CachedMapping map[string]string `json:"mapping"`
	IndexType     string            `json:"index_type"`   // "system" or "user"
	StorageType   string            `json:"storage_type"` // disk, memory, s3
}
