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
	"fmt"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/metadata"
)

// CheckShards if current shard reach the maximum shard size, create a new shard
func (index *Index) CheckShards() error {
	w, err := index.GetWriter()
	if err != nil {
		return err
	}
	_, size := w.DirectoryStats()
	if size > config.Global.Shard.MaxSize {
		return index.NewShard()
	}
	return nil
}

func (index *Index) NewShard() error {
	log.Info().Str("index", index.Name).Int("shard", index.ShardNum).Msg("init new shard")
	// update current shard
	index.UpdateMetadataByShard(index.ShardNum - 1)
	index.lock.Lock()
	index.Shards[index.ShardNum-1].DocTimeMin = index.DocTimeMin
	index.Shards[index.ShardNum-1].DocTimeMax = index.DocTimeMax
	index.DocTimeMin = 0
	index.DocTimeMax = 0
	// create new shard
	index.ShardNum++
	index.Shards = append(index.Shards, &meta.IndexShard{ID: index.ShardNum - 1})
	index.lock.Unlock()
	// store update
	err := metadata.Index.Set(index.Name, index.Index)
	if err != nil {
		return err
	}
	return index.openWriter(index.ShardNum - 1)
}

// GetWriter return the newest shard writer or special shard writer
func (index *Index) GetWriter(shards ...int) (*bluge.Writer, error) {
	var shard int
	if len(shards) == 1 {
		shard = shards[0]
	} else {
		shard = index.ShardNum - 1
	}
	if shard >= index.ShardNum || shard < 0 {
		return nil, errors.New(errors.ErrorTypeRuntimeException, "shard not found")
	}
	index.lock.RLock()
	w := index.Shards[shard].Writer
	index.lock.RUnlock()
	if w != nil {
		return w, nil
	}

	// open writer
	if err := index.openWriter(shard); err != nil {
		return nil, err
	}

	index.lock.RLock()
	w = index.Shards[shard].Writer
	index.lock.RUnlock()

	return w, nil
}

// GetWriters return all shard writers
func (index *Index) GetWriters() ([]*bluge.Writer, error) {
	ws := make([]*bluge.Writer, 0, index.ShardNum)
	for i := index.ShardNum - 1; i >= 0; i-- {
		w, err := index.GetWriter(i)
		if err != nil {
			return nil, err
		}
		ws = append(ws, w)
	}
	return ws, nil
}

// GetReaders return all shard readers
func (index *Index) GetReaders(timeMin, timeMax int64) ([]*bluge.Reader, error) {
	rs := make([]*bluge.Reader, 0, 1)
	chs := make(chan *bluge.Reader, index.ShardNum)
	egLimit := make(chan struct{}, 10)
	eg := errgroup.Group{}
	for i := index.ShardNum - 1; i >= 0; i-- {
		if (timeMin > 0 && index.Shards[i].DocTimeMax > 0 && index.Shards[i].DocTimeMax < timeMin) ||
			(timeMax > 0 && index.Shards[i].DocTimeMin > 0 && index.Shards[i].DocTimeMin > timeMax) {
			continue
		}
		var i = i
		egLimit <- struct{}{}
		eg.Go(func() error {
			defer func() {
				<-egLimit
			}()
			w, err := index.GetWriter(i)
			if err != nil {
				return err
			}
			r, err := w.Reader()
			if err != nil {
				return err
			}
			chs <- r
			return nil
		})
		if index.Shards[i].DocTimeMin < timeMin {
			break
		}
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	close(chs)
	for r := range chs {
		rs = append(rs, r)
	}
	return rs, nil
}

func (index *Index) openWriter(shard int) error {
	var defaultSearchAnalyzer *analysis.Analyzer
	if index.Analyzers != nil {
		defaultSearchAnalyzer = index.Analyzers["default"]
	}

	indexName := fmt.Sprintf("%s/%06x", index.Name, shard)
	var err error
	index.lock.Lock()
	index.Shards[shard].Writer, err = OpenIndexWriter(indexName, index.StorageType, defaultSearchAnalyzer, 0, 0)
	index.lock.Unlock()
	return err
}
