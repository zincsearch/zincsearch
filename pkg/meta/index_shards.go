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
	"sync/atomic"

	"github.com/zinclabs/zinc/pkg/errors"
)

type IndexShards struct {
	ShardNum int64                       `json:"shard_num"`
	Shards   map[string]*IndexFirstShard `json:"shards"`
	lock     sync.RWMutex
}

type IndexFirstShard struct {
	ID       string              `json:"id"`
	ShardNum int64               `json:"shard_num"`
	Shards   []*IndexSecondShard `json:"shards"`
	Stats    *IndexStat          `json:"stats"`
	lock     sync.RWMutex
}

type IndexSecondShard struct {
	ID    int64      `json:"id"`
	Stats *IndexStat `json:"stats"`
}

func NewIndexShards() *IndexShards {
	return &IndexShards{
		Shards: make(map[string]*IndexFirstShard),
	}
}

func (t *IndexShards) Create(id string) (*IndexFirstShard, error) {
	t.lock.RLock()
	_, ok := t.Shards[id]
	t.lock.RUnlock()
	if ok {
		return nil, errors.ErrShardIsExists
	}

	s := &IndexFirstShard{
		ID:       id,
		ShardNum: 1,
		Shards:   []*IndexSecondShard{{ID: 1, Stats: NewIndexStat()}},
		Stats:    NewIndexStat(),
	}
	t.lock.Lock()
	t.Shards[id] = s
	atomic.AddInt64(&t.ShardNum, 1)
	t.lock.Unlock()
	return s, nil
}

func (t *IndexShards) List() []*IndexFirstShard {
	m := make([]*IndexFirstShard, 0)
	t.lock.RLock()
	for _, v := range t.Shards {
		m = append(m, v)
	}
	t.lock.RUnlock()
	return m
}

func (t *IndexShards) Set(shard *IndexFirstShard) error {
	t.lock.Lock()
	_, ok := t.Shards[shard.GetID()]
	if !ok {
		t.Shards[shard.GetID()] = shard
	}
	t.lock.Unlock()
	if ok {
		return errors.ErrShardIsExists
	}
	atomic.AddInt64(&t.ShardNum, 1)
	return nil
}

func (t *IndexShards) Reset(shard *IndexFirstShard) error {
	t.lock.Lock()
	_, ok := t.Shards[shard.GetID()]
	if ok {
		t.Shards[shard.GetID()] = shard
	}
	t.lock.Unlock()
	if !ok {
		return errors.ErrShardNotExists
	}
	return nil
}

func (t *IndexShards) Get(id string) *IndexFirstShard {
	t.lock.RLock()
	s := t.Shards[id]
	t.lock.RUnlock()
	return s
}

func (t *IndexShards) GetShardNum() int64 {
	return atomic.LoadInt64(&t.ShardNum)
}

func (t *IndexShards) Copy() []*IndexFirstShard {
	m := make([]*IndexFirstShard, 0)
	t.lock.RLock()
	for _, v := range t.Shards {
		m = append(m, v.Copy())
	}
	t.lock.RUnlock()
	return m
}

func (t *IndexShards) CopyByID(id string) (*IndexFirstShard, error) {
	t.lock.RLock()
	s, ok := t.Shards[id]
	t.lock.RUnlock()
	if !ok {
		return nil, errors.ErrShardNotExists
	}
	return s.Copy(), nil
}

func (t *IndexFirstShard) GetID() string {
	return t.ID
}

func (t *IndexFirstShard) Create() *IndexSecondShard {
	t.lock.Lock()
	defer t.lock.Unlock()
	s := &IndexSecondShard{ID: t.GetShardNum(), Stats: NewIndexStat()}
	t.Shards = append(t.Shards, s)
	atomic.AddInt64(&t.ShardNum, 1)
	return s
}

func (t *IndexFirstShard) GetShardNum() int64 {
	return atomic.LoadInt64(&t.ShardNum)
}

func (t *IndexFirstShard) List() []*IndexSecondShard {
	m := make([]*IndexSecondShard, 0)
	t.lock.RLock()
	m = append(m, t.Shards...)
	t.lock.RUnlock()
	return m
}

func (t *IndexFirstShard) Copy() *IndexFirstShard {
	s := &IndexFirstShard{
		ID:       t.ID,
		ShardNum: t.GetShardNum(),
		Stats:    t.Stats.Copy(),
	}
	t.lock.RLock()
	for _, v := range t.Shards {
		s.Shards = append(s.Shards, v.Copy())
	}
	t.lock.RUnlock()
	return s
}

func (t *IndexSecondShard) GetID() int64 {
	return t.ID
}

func (t *IndexSecondShard) Copy() *IndexSecondShard {
	return &IndexSecondShard{
		ID:    t.ID,
		Stats: t.Stats.Copy(),
	}
}
