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

	"github.com/jinzhu/copier"
)

type IndexSettings struct {
	NumberOfShards   int64          `json:"number_of_shards,omitempty"`
	NumberOfReplicas int64          `json:"number_of_replicas,omitempty"`
	Analysis         *IndexAnalysis `json:"analysis,omitempty"`
	lock             sync.RWMutex
}

type IndexAnalysis struct {
	Analyzer    map[string]*Analyzer   `json:"analyzer,omitempty"`
	CharFilter  map[string]interface{} `json:"char_filter,omitempty"`
	Tokenizer   map[string]interface{} `json:"tokenizer,omitempty"`
	TokenFilter map[string]interface{} `json:"token_filter,omitempty"`
	Filter      map[string]interface{} `json:"filter,omitempty"` // compatibility with es, alias for TokenFilter
}

func NewIndexSettings() *IndexSettings {
	return &IndexSettings{
		Analysis: &IndexAnalysis{},
	}
}

func (t *IndexSettings) Set(sets *IndexSettings) {
	if sets.NumberOfShards > 0 {
		t.SetShards(sets.NumberOfShards)
	}
	if sets.NumberOfReplicas > 0 {
		t.SetReplicas(sets.NumberOfReplicas)
	}
	if sets.Analysis != nil {
		t.SetAnalysis(sets.Analysis)
	}
}

func (t *IndexSettings) GetShards() int64 {
	return atomic.LoadInt64(&t.NumberOfShards)
}

func (t *IndexSettings) SetShards(num int64) {
	atomic.StoreInt64(&t.NumberOfShards, num)
}

func (t *IndexSettings) GetReplicas() int64 {
	return atomic.LoadInt64(&t.NumberOfReplicas)
}

func (t *IndexSettings) SetReplicas(num int64) {
	atomic.StoreInt64(&t.NumberOfReplicas, num)
}

func (t *IndexSettings) GetAnalysis() *IndexAnalysis {
	t.lock.RLock()
	a := t.Analysis
	t.lock.RUnlock()
	return a
}

func (t *IndexSettings) SetAnalysis(analysis *IndexAnalysis) {
	t.lock.Lock()
	t.Analysis = analysis
	t.lock.Unlock()
}

func (t *IndexSettings) Copy() *IndexSettings {
	ana := new(IndexAnalysis)
	t.lock.RLock()
	copier.Copy(ana, t.Analysis)
	t.lock.RUnlock()
	return &IndexSettings{
		NumberOfShards:   t.NumberOfShards,
		NumberOfReplicas: t.NumberOfReplicas,
		Analysis:         ana,
	}
}
