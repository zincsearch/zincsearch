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
	"strings"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/uquery"
	"github.com/zinclabs/zinc/pkg/uquery/timerange"
)

func MultiSearch(indexNames []string, query *meta.ZincQuery) (*meta.SearchResponse, error) {
	var mappings *meta.Mappings
	var analyzers map[string]*analysis.Analyzer
	var readers []*bluge.Reader
	var shardNum int
	indexMap := make(map[string]struct{})

	timeMin, timeMax := timerange.Query(query.Query)
	for _, index := range ZINC_INDEX_LIST.List() {
		for _, indexName := range indexNames {
			if _, ok := indexMap[index.Name]; ok {
				continue
			}
			if indexName == "" || (indexName != "" && strings.HasPrefix(index.Name, indexName[:len(indexName)-1])) {
				reader, err := index.GetReaders(timeMin, timeMax)
				if err != nil {
					return nil, err
				}
				readers = append(readers, reader...)
				shardNum += index.ShardNum
				if mappings == nil {
					mappings = index.Mappings
					analyzers = index.Analyzers
				}
				indexMap[index.Name] = struct{}{}
			}
		}
	}
	defer func() {
		for _, reader := range readers {
			reader.Close()
		}
	}()

	if len(readers) == 0 {
		return nil, fmt.Errorf("core.MultiSearchV2: error accessing reader: no index found")
	}

	searchRequest, err := uquery.ParseQueryDSL(query, mappings, analyzers)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var cancel context.CancelFunc
	if query.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(query.Timeout)*time.Second)
		defer cancel()
	}

	dmi, err := bluge.MultiSearch(ctx, searchRequest, readers...)
	if err != nil {
		log.Printf("core.MultiSearchV2: error executing search: %s", err.Error())
		if err == context.DeadlineExceeded {
			return &meta.SearchResponse{
				TimedOut: true,
				Error:    err.Error(),
				Hits:     meta.Hits{Hits: []meta.Hit{}},
			}, nil
		}
		return nil, err
	}

	return searchV2(shardNum, len(readers), dmi, query, mappings)
}
