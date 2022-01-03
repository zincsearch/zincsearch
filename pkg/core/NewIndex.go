package core

import (
	"github.com/blugelabs/bluge"
	"github.com/prabhatsharma/zinc/directory"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

// NewIndex creates an instance of a physical zinc index that can be used to store and retrieve data.
func NewIndex(name string, storageType string) (*Index, error) {

	var config bluge.Config

	if storageType == "disk" {
		DATA_PATH := zutils.GetEnv("DATA_PATH", "./data")

		config = bluge.DefaultConfig(DATA_PATH + "/" + name)
	} else if storageType == "s3" {
		S3_BUCKET := zutils.GetEnv("S3_BUCKET", "zinc1")
		config = directory.GetS3Config(S3_BUCKET, name)
	}

	writer, err := bluge.OpenWriter(config)

	if err != nil {
		return nil, err
	}

	index := &Index{
		Name:   name,
		Writer: writer,
		StorageType: storageType,
	}

	mapping, err := index.GetStoredMapping()
	if err != nil {
		return nil, err
	}

	index.CachedMapping = mapping

	return index, nil
}

func IndexExists(index string) (bool, string) {
	if _, ok := ZINC_INDEX_LIST[index]; ok {
		return true, ZINC_INDEX_LIST[index].StorageType
	}

	return false, ""
}
