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
	"context"
	"fmt"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"

	zincanalysis "github.com/zinclabs/zinc/pkg/uquery/analysis"
)

var systemIndexList = []string{"_index_mapping", "_index_template", "_index", "_metadata", "_users"}

func LoadZincSystemIndexes() (map[string]*Index, error) {
	indexList := make(map[string]*Index)
	for _, index := range systemIndexList {
		log.Info().Msgf("Loading system index... [%s:%s]", index, "disk")
		writer, err := LoadIndexWriter(index, "disk", nil)
		if err != nil {
			return nil, err
		}
		indexList[index] = &Index{
			Name:        index,
			IndexType:   "system",
			StorageType: "disk",
			Writer:      writer,
		}
	}

	return indexList, nil
}

func LoadZincIndexesFromMeta() (map[string]*Index, error) {
	query := bluge.NewMatchAllQuery()
	searchRequest := bluge.NewAllMatches(query).WithStandardAggregations()
	reader, _ := ZINC_SYSTEM_INDEX_LIST["_index"].Writer.Reader()
	defer reader.Close()

	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		return nil, fmt.Errorf("core.LoadZincIndexesFromMeta: error executing search: %s", err.Error())
	}

	indexList := make(map[string]*Index)
	next, err := dmi.Next()
	for err == nil && next != nil {
		index := &Index{IndexType: "user", StorageType: "disk"}
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			switch field {
			case "name":
				index.Name = string(value)
			case "index_type":
				index.IndexType = string(value)
			case "storage_type":
				index.StorageType = string(value)
			case "settings":
				_ = json.Unmarshal(value, &index.Settings)
			case "mappings":
				_ = json.Unmarshal(value, &index.CachedMappings)
			default:
			}
			return true
		})

		log.Info().Msgf("Loading user   index... [%s:%s]", index.Name, index.StorageType)
		if err != nil {
			log.Printf("core.LoadZincIndexesFromMeta: error accessing stored fields: %s", err.Error())
		}

		// load index analysis
		if index.Settings != nil && index.Settings.Analysis != nil {
			index.CachedAnalyzers, err = zincanalysis.RequestAnalyzer(index.Settings.Analysis)
			if err != nil {
				log.Printf("core.LoadZincIndexesFromMeta: error parse stored analysis: %s", err.Error())
			}
		}

		// load index data
		var defaultSearchAnalyzer *analysis.Analyzer
		if index.CachedAnalyzers != nil {
			defaultSearchAnalyzer = index.CachedAnalyzers["default"]
		}
		index.Writer, err = LoadIndexWriter(index.Name, index.StorageType, defaultSearchAnalyzer)
		if err != nil {
			log.Error().Msgf("Loading user   index... [%s:%s] index writer error: %s", index.Name, index.StorageType, err.Error())
		}

		// load index docs count
		index.DocsCount, _ = index.LoadDocsCount()

		// load index size
		index.ReLoadStorageSize()

		indexList[index.Name] = index

		next, err = dmi.Next()
	}

	return indexList, nil
}

func CloseIndexes() {
	for _, index := range ZINC_INDEX_LIST {
		_ = index.Close()
	}
	for _, index := range ZINC_SYSTEM_INDEX_LIST {
		_ = index.Close()
	}
}
