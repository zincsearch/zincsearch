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
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/meta"
)

var ZINC_INDEX_LIST map[string]*Index
var ZINC_SYSTEM_INDEX_LIST map[string]*Index

func init() {
	var err error
	ZINC_SYSTEM_INDEX_LIST, err = LoadZincSystemIndexes()
	if err != nil {
		log.Fatal().Msgf("Error loading system index: %s", err.Error())
	}

	ZINC_INDEX_LIST, _ = LoadZincIndexesFromMeta()
	if err != nil {
		log.Error().Msgf("Error loading user index: %s", err.Error())
	}
}

type Index struct {
	Name                string                        `json:"name"`
	IndexType           string                        `json:"index_type"`   // "system" or "user"
	StorageType         string                        `json:"storage_type"` // disk, memory, s3
	DocsCount           int64                         `json:"docs_count"`   // cached docs count of the index
	StorageSize         float64                       `json:"size"`         // cached size of the index
	StorageSizeNextTime time.Time                     `json:"-"`            // control the update of the storage size
	Mappings            map[string]interface{}        `json:"mappings"`
	Settings            *meta.IndexSettings           `json:"settings"`
	CachedAnalyzers     map[string]*analysis.Analyzer `json:"-"`
	CachedMappings      *meta.Mappings                `json:"-"`
	Writer              *bluge.Writer                 `json:"-"`
}

type IndexTemplate struct {
	Name          string         `json:"name"`
	Timestamp     time.Time      `json:"@timestamp"`
	IndexTemplate *meta.Template `json:"index_template"`
}
