package core

import (
	"github.com/blugelabs/bluge"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

type Index struct {
	Name           string         `json:"name"`
	IndexType      string         `json:"index_type"`   // "system" or "user"
	StorageType    string         `json:"storage_type"` // disk, memory, s3
	Mappings       *meta.Mappings `json:"mappings"`
	CachedMappings *meta.Mappings `json:"-"`
	Writer         *bluge.Writer  `json:"-"`
}

var ZINC_INDEX_LIST map[string]*Index
var ZINC_SYSTEM_INDEX_LIST map[string]*Index

func init() {
	ZINC_SYSTEM_INDEX_LIST, _ = LoadZincSystemIndexes()
	ZINC_INDEX_LIST, _ = LoadZincIndexesFromDisk()
	s3List, _ := LoadZincIndexesFromS3()
	// Load the indexes from disk.
	for k, v := range s3List {
		ZINC_INDEX_LIST[k] = v
	}

	minioList, _ := LoadZincIndexesFromMinIO()
	// Load the indexes from disk.
	for k, v := range minioList {
		ZINC_INDEX_LIST[k] = v
	}
}
