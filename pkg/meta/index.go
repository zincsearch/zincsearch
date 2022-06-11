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

package meta

import (
	"time"

	"github.com/blugelabs/bluge"
)

type Index struct {
	Name        string         `json:"name"`
	StorageType string         `json:"storage_type"`
	StorageSize uint64         `json:"storage_size"`
	DocNum      uint64         `json:"doc_num"`
	DocTimeMin  int64          `json:"doc_time_min"`
	DocTimeMax  int64          `json:"doc_time_max"`
	ShardNum    int            `json:"shard_num"`
	Shards      []*IndexShard  `json:"shards"`
	Settings    *IndexSettings `json:"settings,omitempty"`
	Mappings    *Mappings      `json:"mappings,omitempty"`
	CreateAt    time.Time      `json:"create_at"`
	UpdateAt    time.Time      `json:"update_at"`
}

type IndexShard struct {
	ID          int           `json:"id"`
	DocTimeMin  int64         `json:"doc_time_min"`
	DocTimeMax  int64         `json:"doc_time_max"`
	DocNum      uint64        `json:"doc_num"`
	StorageSize uint64        `json:"storage_size"`
	Writer      *bluge.Writer `json:"-"`
}

type IndexSimple struct {
	Name        string                 `json:"name"`
	StorageType string                 `json:"storage_type"`
	Settings    *IndexSettings         `json:"settings,omitempty"`
	Mappings    map[string]interface{} `json:"mappings,omitempty"`
}

type IndexSettings struct {
	NumberOfShards   int            `json:"number_of_shards,omitempty"`
	NumberOfReplicas int            `json:"number_of_replicas,omitempty"`
	Analysis         *IndexAnalysis `json:"analysis,omitempty"`
}

type IndexAnalysis struct {
	Analyzer    map[string]*Analyzer   `json:"analyzer,omitempty"`
	CharFilter  map[string]interface{} `json:"char_filter,omitempty"`
	Tokenizer   map[string]interface{} `json:"tokenizer,omitempty"`
	TokenFilter map[string]interface{} `json:"token_filter,omitempty"`
	Filter      map[string]interface{} `json:"filter,omitempty"` // compatibility with es, alias for TokenFilter
}

type SortIndex []*Index

func (t SortIndex) Len() int {
	return len(t)
}

func (t SortIndex) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t SortIndex) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}
