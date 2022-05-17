/* Copyright 2022 Zinc Labs Inc. and Contributors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/goccy/go-json"

	"github.com/zinclabs/zinc/pkg/bluge/directory"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/zutils"
)

// NewIndex creates an instance of a physical zinc index that can be used to store and retrieve data.
func NewIndex(name, storageType string, defaultSearchAnalyzer *analysis.Analyzer) (*Index, error) {
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
		storageType = "disk"
		dataPath = zutils.GetEnv("ZINC_DATA_PATH", "./data")
		config = bluge.DefaultConfig(dataPath + "/" + name)
	}

	if defaultSearchAnalyzer != nil {
		config.DefaultSearchAnalyzer = defaultSearchAnalyzer
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

	return index, nil
}

// LoadIndexWriter load the index writer from the storage
func LoadIndexWriter(name string, storageType string, defaultSearchAnalyzer *analysis.Analyzer) (*bluge.Writer, error) {
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

	if defaultSearchAnalyzer != nil {
		config.DefaultSearchAnalyzer = defaultSearchAnalyzer
	}

	return bluge.OpenWriter(config)
}

// storeIndex stores the index to metadata
func StoreIndex(index *Index) error {
	if index.Settings == nil {
		index.Settings = new(meta.IndexSettings)
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
		return fmt.Errorf("core.StoreIndex: index: %s, error: %s", index.Name, err.Error())
	}

	// cache index
	ZINC_INDEX_LIST[index.Name] = index

	return nil
}

func GetIndex(name string) (*Index, bool) {
	index, ok := ZINC_INDEX_LIST[name]
	return index, ok
}
