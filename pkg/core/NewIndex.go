package core

import (
	"os"

	"github.com/blugelabs/bluge"
)

func NewIndex(name string) (*Index, error) {
	DATA_PATH := ""
	if os.Getenv("DATA_PATH") == "" {
		DATA_PATH = "./data"
		// DATA_PATH = "/Users/prabhat/projects/prabhatsharma/zinc/data"
	} else {
		DATA_PATH = os.Getenv("DATA_PATH")
	}

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
