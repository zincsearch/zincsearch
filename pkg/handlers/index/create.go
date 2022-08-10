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

package index

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/meta"
	zincanalysis "github.com/zinclabs/zinc/pkg/uquery/analysis"
	"github.com/zinclabs/zinc/pkg/uquery/mappings"
	"github.com/zinclabs/zinc/pkg/zutils"
)

// @Id CreateIndex
// @Summary Create index
// @Tags    Index
// @Accept  json
// @Produce json
// @Param   data body meta.IndexSimple true "Index data"
// @Success 200 {object} meta.HTTPResponseIndex
// @Failure 400 {object} meta.HTTPResponseError
// @Router /api/index [post]
func Create(c *gin.Context) {
	var newIndex meta.IndexSimple
	if err := zutils.GinBindJSON(c, &newIndex); err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	indexName := c.Param("target")
	err := CreateIndexWorker(&newIndex, indexName)
	if err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, meta.HTTPResponseIndex{
		Message:     "ok",
		Index:       newIndex.Name,
		StorageType: newIndex.StorageType,
	})
}

// @Id ESCreateIndex
// @Summary Create index for compatible ES
// @Tags    Index
// @Produce json
// @Param   index path  string  true  "Index"
// @Param   data  body  meta.IndexSimple true "Index data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} meta.HTTPResponse
// @Router /es/{index} [put]
func CreateES(c *gin.Context) {
	indexName := c.Param("target")

	var newIndex meta.IndexSimple
	if err := zutils.GinBindJSON(c, &newIndex); err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	// TODO: migrate ES settings to Zinc specific
	newIndex.Settings = nil

	// default the storage_type to disk, to provide the best possible integration
	newIndex.StorageType = "disk"

	err := CreateIndexWorker(&newIndex, indexName)
	if err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"acknowledged":        true,
		"shards_acknowledged": true,
		"index":               newIndex.Name,
	})
}

func CreateIndexWorker(newIndex *meta.IndexSimple, indexName string) error {
	if newIndex.Name == "" && indexName != "" {
		newIndex.Name = indexName
	}

	if newIndex.Name == "" {
		return errors.ErrIndexIsEmpty
	}

	if _, ok := core.GetIndex(newIndex.Name); ok {
		return errors.ErrIndexIsExists
	}

	if newIndex.Settings == nil {
		newIndex.Settings = meta.NewIndexSettings()
	}
	analyzers, err := zincanalysis.RequestAnalyzer(newIndex.Settings.Analysis)
	if err != nil {
		return err
	}

	mappings, err := mappings.Request(analyzers, newIndex.Mappings)
	if err != nil {
		return err
	}

	if newIndex.Settings != nil && newIndex.Settings.NumberOfShards != 0 {
		if newIndex.ShardNum == 0 {
			newIndex.ShardNum = newIndex.Settings.NumberOfShards
		}
	}
	if newIndex.ShardNum == 0 {
		newIndex.ShardNum = config.Global.Shard.Num
	}
	index, exists, err := core.GetOrCreateIndex(newIndex.Name, newIndex.StorageType, newIndex.ShardNum)
	if err != nil {
		return err
	}
	if exists {
		return errors.ErrIndexIsExists
	}

	// update settings
	if err := index.SetSettings(newIndex.Settings, true); err != nil {
		return err
	}

	// update mappings
	if err := index.SetMappings(mappings, true); err != nil {
		return err
	}

	// update analyzers
	index.SetAnalyzers(analyzers)

	return nil
}
