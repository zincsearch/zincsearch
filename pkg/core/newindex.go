package core

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/blugelabs/bluge"

	"github.com/prabhatsharma/zinc/pkg/directory"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

// NewIndex creates an instance of a physical zinc index that can be used to store and retrieve data.
func NewIndex(name string, storageType string, useNewIndexMeta int) (*Index, error) {
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

	// use template
	if err = index.UseTemplate(); err != nil {
		return nil, err
	}

	if useNewIndexMeta == NotCompatibleNewIndexMeta {
		mappings, err := index.GetStoredMapping()
		if err != nil {
			return nil, err
		}

		if mappings != nil && len(mappings.Properties) > 0 {
			index.CachedMappings = mappings
		}
	}

	return index, nil
}

// LoadIndexWriter load the index writer from the storage
func LoadIndexWriter(name string, storageType string) (*bluge.Writer, error) {
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

	return bluge.OpenWriter(config)
}

// storeIndex stores the index to metadata
func StoreIndex(index *Index, needUpdate bool) error {
	bdoc := bluge.NewDocument(index.Name)
	bdoc.AddField(bluge.NewKeywordField("name", index.Name).StoreValue().Sortable())
	bdoc.AddField(bluge.NewKeywordField("index_type", index.IndexType).StoreValue().Sortable())
	bdoc.AddField(bluge.NewKeywordField("storage_type", index.StorageType).StoreValue().Sortable())

	settingByteVal, _ := json.Marshal(&index.Settings)
	bdoc.AddField(bluge.NewStoredOnlyField("settings", settingByteVal))
	mappingByteVal, _ := json.Marshal(&index.CachedMappings)
	bdoc.AddField(bluge.NewStoredOnlyField("mappings", mappingByteVal))

	bdoc.AddField(bluge.NewDateTimeField("@timestamp", time.Now()).StoreValue().Sortable().Aggregatable())
	bdoc.AddField(bluge.NewStoredOnlyField("_source", nil))
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil)) // Add _all field that can be used for search

	var err error
	indexWriter := ZINC_SYSTEM_INDEX_LIST["_index"].Writer
	if needUpdate {
		err = indexWriter.Update(bdoc.ID(), bdoc)
	} else {
		err = indexWriter.Insert(bdoc)
	}
	if err != nil {
		return fmt.Errorf("core.StoreIndex: error updating document: %v", err)
	}

	// cache index
	ZINC_INDEX_LIST[index.Name] = index

	return nil
}

func GetIndex(indexName string) (*Index, bool) {
	index, ok := ZINC_INDEX_LIST[indexName]
	return index, ok
}
