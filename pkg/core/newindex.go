package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"

	"github.com/prabhatsharma/zinc/pkg/directory"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

// NewIndex creates an instance of a physical zinc index that can be used to store and retrieve data.
func NewIndex(name string, storageType string, useNewIndexMeta int) (*Index, error) {
	if name == "" {
		return nil, fmt.Errorf("core.NewIndex: index name cannot be empty")
	}
	if strings.HasPrefix(name, "_") {
		return nil, fmt.Errorf("core.NewIndex: index name cannot start with _")
	}

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
func StoreIndex(index *Index) error {
	if index.Settings == nil {
		index.Settings = meta.NewIndexSettings()
	}
	if index.CachedAnalyzers == nil {
		index.CachedAnalyzers = make(map[string]*analysis.Analyzer)
	}
	if index.CachedMappings == nil {
		index.CachedMappings = meta.NewMappings()
	}

	bdoc := bluge.NewDocument(index.Name)
	bdoc.AddField(bluge.NewKeywordField("name", index.Name).StoreValue().Sortable())
	bdoc.AddField(bluge.NewKeywordField("index_type", index.IndexType).StoreValue().Sortable())
	bdoc.AddField(bluge.NewKeywordField("storage_type", index.StorageType).StoreValue().Sortable())

	settingByteVal, _ := json.Marshal(&index.Settings)
	bdoc.AddField(bluge.NewStoredOnlyField("settings", settingByteVal))
	mappingsByteVal, _ := json.Marshal(&index.CachedMappings)
	bdoc.AddField(bluge.NewStoredOnlyField("mappings", mappingsByteVal))

	bdoc.AddField(bluge.NewDateTimeField("@timestamp", time.Now()).StoreValue().Sortable().Aggregatable())
	bdoc.AddField(bluge.NewStoredOnlyField("_source", nil))
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil)) // Add _all field that can be used for search

	err := ZINC_SYSTEM_INDEX_LIST["_index"].Writer.Update(bdoc.ID(), bdoc)
	if err != nil {
		return fmt.Errorf("core.StoreIndex: index: %s, error: %v", index.Name, err)
	}

	// cache index
	ZINC_INDEX_LIST[index.Name] = index

	return nil
}

func DeleteIndex(name string) error {
	bdoc := bluge.NewDocument(name)
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil))
	err := ZINC_SYSTEM_INDEX_LIST["_index"].Writer.Delete(bdoc.ID())
	if err != nil {
		return fmt.Errorf("core.DeleteIndex: error deleting template: %v", err)
	}

	return nil
}

func GetIndex(name string) (*Index, bool) {
	index, ok := ZINC_INDEX_LIST[name]
	return index, ok
}
