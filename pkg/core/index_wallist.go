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
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/zinclabs/zinc/pkg/config"
)

// Record opened WAL used to do consume
var ZINC_INDEX_WAL_LIST IndexWALList

type IndexWALList struct {
	Indexes map[string]*Index
	lock    sync.RWMutex
}

func init() {
	ZINC_INDEX_WAL_LIST.Indexes = make(map[string]*Index)
	go ZINC_INDEX_WAL_LIST.ConsumeWAL()
}

func (t *IndexWALList) Add(index *Index) {
	t.lock.Lock()
	t.Indexes[index.GetName()] = index
	t.lock.Unlock()
}

func (t *IndexWALList) Remove(name string) {
	t.lock.Lock()
	delete(t.Indexes, name)
	t.lock.Unlock()
}

func (t *IndexWALList) Len() int {
	t.lock.RLock()
	n := len(t.Indexes)
	t.lock.RUnlock()
	return n
}

func (t *IndexWALList) ConsumeWAL() {
	interval, err := parseInterval(config.Global.WalSyncInterval)
	if err != nil {
		log.Fatal().Err(err).Msg("consume ParseInterval")
	}
	tick := time.NewTicker(interval)
	for range tick.C {
		eg := &errgroup.Group{}
		eg.SetLimit(config.Global.ReadGorutineNum)
		indexClosed := make(chan string, t.Len())
		t.lock.RLock()
		for _, idx := range t.Indexes {
			var index = idx
			eg.Go(func() error {
				select {
				case <-index.close:
					indexClosed <- index.GetName()
					return nil
				default:
					// continue
				}
				index.ConsumeWAL()
				return nil
			})
		}
		_ = eg.Wait()
		t.lock.RUnlock()
		close(indexClosed)

		// check index closed
		for name := range indexClosed {
			t.Remove(name)
		}
	}
}
