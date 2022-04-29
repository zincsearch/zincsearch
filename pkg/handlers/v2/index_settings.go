// Copyright 2022 Zinc Labs Inc. and Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v2

import (
	"net/http"

	"github.com/blugelabs/bluge/analysis"
	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zinc/pkg/core"
	meta "github.com/zinclabs/zinc/pkg/meta/v2"
	zincanalysis "github.com/zinclabs/zinc/pkg/uquery/v2/analysis"
)

func GetIndexSettings(c *gin.Context) {
	indexName := c.Param("target")
	index, exists := core.GetIndex(indexName)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "index " + indexName + " does not exists"})
		return
	}

	settings := index.Settings
	if settings == nil {
		settings = new(meta.IndexSettings)
	}

	c.JSON(http.StatusOK, gin.H{index.Name: gin.H{"settings": settings}})
}

func UpdateIndexSettings(c *gin.Context) {
	indexName := c.Param("target")
	if indexName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index.name should be not empty"})
		return
	}

	var newIndex core.Index
	if err := c.BindJSON(&newIndex); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newIndex.Settings == nil {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
		return
	}

	analyzers, err := zincanalysis.RequestAnalyzer(newIndex.Settings.Analysis)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	index, exists := core.GetIndex(indexName)
	if exists {
		// it can only change settings.NumberOfReplicas when index exists
		if newIndex.Settings.NumberOfReplicas > 0 {
			index.Settings.NumberOfReplicas = newIndex.Settings.NumberOfReplicas
		}
		if newIndex.Settings.Analysis != nil && len(newIndex.Settings.Analysis.Analyzer) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't update analyzer for existing index"})
			return
		}
		// store index
		core.StoreIndex(index)

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
		return
	}

	var defaultSearchAnalyzer *analysis.Analyzer
	if analyzers != nil {
		defaultSearchAnalyzer = analyzers["default"]
	}
	index, err = core.NewIndex(indexName, newIndex.StorageType, core.UseNewIndexMeta, defaultSearchAnalyzer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// update settings
	index.SetSettings(newIndex.Settings)

	// update analyzers
	index.SetAnalyzers(analyzers)

	// store index
	core.StoreIndex(index)

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
