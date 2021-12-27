package core

import (
	"github.com/blugelabs/bluge"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

// NewIndex creates an instance of a physical zinc index that can be used to store and retrieve data.
func NewIndex(name string) (*Index, error) {
	DATA_PATH := zutils.GetEnv("DATA_PATH", "./data")

	config := bluge.DefaultConfig(DATA_PATH + "/" + name)

	writer, err := bluge.OpenWriter(config)

	if err != nil {
		return nil, err
	}

	index := &Index{
		Name:   name,
		Writer: writer,
	}

	mapping, err := index.GetMappingFromDisk()
	if err != nil {
		return nil, err
	}

	index.CachedMapping = mapping

	return index, nil
}

func IndexExists(index string) bool {
	if _, ok := ZINC_INDEX_LIST[index]; ok {
		return true
	}

	return false
}
