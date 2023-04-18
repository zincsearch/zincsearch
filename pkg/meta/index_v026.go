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
	"sync"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/zincsearch/zincsearch/pkg/wal"
)

// IndexV026 compatible for v0.2.6 index
type IndexV026 struct {
	Name        string            `json:"name"`
	StorageType string            `json:"storage_type"`
	StorageSize uint64            `json:"storage_size"`
	DocNum      uint64            `json:"doc_num"`
	DocTimeMin  int64             `json:"doc_time_min"`
	DocTimeMax  int64             `json:"doc_time_max"`
	ShardNum    int64             `json:"shard_num"`
	Shards      []*IndexShardV026 `json:"shards"`
	WAL         *wal.Log          `json:"-"`
	WALSize     uint64            `json:"wal_size"`
	Settings    *IndexSettings    `json:"settings,omitempty"`
	Mappings    *Mappings         `json:"mappings,omitempty"`
	CreateAt    time.Time         `json:"create_at"`
	UpdateAt    time.Time         `json:"update_at"`
}

type IndexShardV026 struct {
	ID          int64         `json:"id"`
	DocTimeMin  int64         `json:"doc_time_min"`
	DocTimeMax  int64         `json:"doc_time_max"`
	DocNum      uint64        `json:"doc_num"`
	StorageSize uint64        `json:"storage_size"`
	Writer      *bluge.Writer `json:"-"`
	Lock        sync.RWMutex  `json:"-"`
}
