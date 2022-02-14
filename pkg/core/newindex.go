package core

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"

	"github.com/prabhatsharma/zinc/pkg/directory"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
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
	mappings, err := index.GetStoredMappings()
	if err != nil {
		return nil, err
	}

	index.CachedMappings = mappings

	return index, nil
}

func GetIndex(indexName string) (*Index, bool) {
	index, ok := ZINC_INDEX_LIST[indexName]
	return index, ok
}

func FormatMappings(mappings meta.Mappings) (*meta.Mappings, error) {
	// copy a mappings
	newmappings := new(meta.Mappings)
	zutils.StructToStruct(mappings, newmappings)

	// format mappings
	for field, prop := range newmappings.Properties {
		prop.Type = strings.ToLower(prop.Type)
		switch prop.Type {
		case "text", "keyword", "numeric", "bool", "time":
			continue // ptype can be used as is
		case "integer", "double", "long":
			prop.Type = "numeric"
		case "boolean":
			prop.Type = "bool"
		case "date", "datetime":
			prop.Type = "time"
		default:
			return nil, fmt.Errorf("[mappings] doesn't type: [%s] for field [%s]", prop.Type, field)
		}
	}

	return newmappings, nil
}
