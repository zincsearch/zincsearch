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

package handlers

import (
	"errors"
	"net/http"

	"github.com/blugelabs/bluge/analysis"
	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zinc/pkg/core"
	meta "github.com/zinclabs/zinc/pkg/meta/v2"
	zincanalysis "github.com/zinclabs/zinc/pkg/uquery/v2/analysis"
	"github.com/zinclabs/zinc/pkg/uquery/v2/mappings"
)

func CreateIndex(c *gin.Context) {
	var newIndex core.Index
	if err := c.BindJSON(&newIndex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	indexName := c.Param("target")

	err := CreateIndexWorker(&newIndex, indexName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "index created",
		"index":        newIndex.Name,
		"storage_type": newIndex.StorageType,
	})
}

func CreateIndexWorker(newIndex *core.Index, indexName string) error {

	if newIndex.Name == "" && indexName != "" {
		newIndex.Name = indexName
	}

	if newIndex.Name == "" {
		return errors.New("index.name should be not empty")
	}

	if _, ok := core.GetIndex(newIndex.Name); ok {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "index [" + newIndex.Name + "] already exists"})
		return errors.New("index [" + newIndex.Name + "] already exists")
	}

	if newIndex.Settings == nil {
		newIndex.Settings = new(meta.IndexSettings)
	}
	analyzers, err := zincanalysis.RequestAnalyzer(newIndex.Settings.Analysis)
	if err != nil {
		return errors.New(err.Error())
	}

	mappings, err := mappings.Request(analyzers, newIndex.Mappings)
	if err != nil {
		return errors.New(err.Error())
	}

	var defaultSearchAnalyzer *analysis.Analyzer
	if analyzers != nil {
		defaultSearchAnalyzer = analyzers["default"]
	}
	index, err := core.NewIndex(newIndex.Name, newIndex.StorageType, core.UseNewIndexMeta, defaultSearchAnalyzer)
	if err != nil {
		return errors.New(err.Error())
	}

	// update settings
	index.SetSettings(newIndex.Settings)

	// update analyzers
	index.SetAnalyzers(analyzers)

	// update mappings
	index.SetMappings(mappings)

	// store index
	err = core.StoreIndex(index)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
