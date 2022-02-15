package core

import (
	"github.com/blugelabs/bluge"

	"github.com/prabhatsharma/zinc/pkg/directory"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

// NewIndex creates an instance of a physical zinc index that can be used to store and retrieve data.
func NewIndex(name string, storageType string) (*Index, error) {
	var dataPath string
	var config bluge.Config
	switch storageType {
	case "s3":
		dataPath = zutils.GetEnv("ZINC_S3_BUCKET", "")
		config = directory.GetS3Config(dataPath, name)
	case "minio":
		dataPath = zutils.GetEnv("ZINC_MINIO_BUCKET", "")
		config = directory.GetMinIOConfig(dataPath, name)
	default:
		dataPath = zutils.GetEnv("ZINC_DATA_PATH", "./data")
		config = bluge.DefaultConfig(dataPath + "/" + name)
	}

	writer, err := bluge.OpenWriter(config)
	if err != nil {
		return nil, err
	}

	index := &Index{
		Name:        name,
		Writer:      writer,
		StorageType: storageType,
	}
	mappings, err := index.GetStoredMapping()
	if err != nil {
		return nil, err
	}

	if mappings != nil && len(mappings.Properties) > 0 {
		index.CachedMappings = mappings
	}

	return index, nil
}

func GetIndex(indexName string) (*Index, bool) {
	index, ok := ZINC_INDEX_LIST[indexName]
	return index, ok
}
