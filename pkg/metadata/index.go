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

package metadata

import (
	"strings"

	"github.com/goccy/go-json"
	"golang.org/x/sync/errgroup"

	"github.com/zinclabs/zinc/pkg/meta"
)

type index struct{}

var Index = new(index)

func (t *index) Get(id string) (*meta.Index, error) {
	idx := meta.NewIndex("", "", "")
	eg := errgroup.Group{}
	eg.Go(func() error {
		data, err := db.Get(t.key(id, "meta"))
		if err != nil {
			return err
		}
		return json.Unmarshal(data, idx.Meta)
	})
	eg.Go(func() error {
		data, err := db.Get(t.key(id, "stats"))
		if err != nil {
			return err
		}
		return json.Unmarshal(data, idx.Stats)
	})
	eg.Go(func() error {
		data, err := db.Get(t.key(id, "settings"))
		if err != nil {
			return err
		}
		idx.Settings = new(meta.IndexSettings)
		return json.Unmarshal(data, idx.Settings)
	})
	eg.Go(func() error {
		data, err := db.Get(t.key(id, "mappings"))
		if err != nil {
			return err
		}
		idx.Mappings = meta.NewMappings()
		return json.Unmarshal(data, idx.Mappings)
	})
	eg.Go(func() error {
		shards, err := t.GetShards(id)
		if err != nil {
			return err
		}
		for _, shard := range shards.Shards {
			if err := idx.Shards.Set(shard); err != nil {
				return err
			}
		}
		return nil
	})
	err := eg.Wait()

	return idx, err
}

func (t *index) GetMeta(id string) (*meta.IndexMeta, error) {
	data, err := db.Get(t.key(id, "meta"))
	if err != nil {
		return nil, err
	}
	val := meta.NewIndexMeta()
	err = json.Unmarshal(data, val)
	return val, err
}

func (t *index) GetStats(id string) (*meta.IndexStat, error) {
	data, err := db.Get(t.key(id, "stats"))
	if err != nil {
		return nil, err
	}
	val := meta.NewIndexStat()
	err = json.Unmarshal(data, val)
	return val, err
}

func (t *index) GetSettings(id string) (*meta.IndexSettings, error) {
	data, err := db.Get(t.key(id, "settings"))
	if err != nil {
		return nil, err
	}
	val := meta.NewIndexSettings()
	err = json.Unmarshal(data, val)
	return val, err
}

func (t *index) GetMappings(id string) (*meta.Mappings, error) {
	data, err := db.Get(t.key(id, "mappings"))
	if err != nil {
		return nil, err
	}
	val := meta.NewMappings()
	err = json.Unmarshal(data, val)
	return val, err
}

func (t *index) GetShards(id string) (*meta.IndexShards, error) {
	data, err := db.List(t.key(id, "shards/"), 0, 0)
	if err != nil {
		return nil, err
	}
	shards := meta.NewIndexShards()
	for _, d := range data {
		shard := new(meta.IndexFirstShard)
		if err = json.Unmarshal(d, shard); err != nil {
			return nil, err
		}
		if err = shards.Set(shard); err != nil {
			return nil, err
		}
	}
	return shards, err
}

func (t *index) GetShard(id, shard string) (*meta.IndexFirstShard, error) {
	data, err := db.Get(t.key(id, "shards", shard))
	if err != nil {
		return nil, err
	}
	val := new(meta.IndexFirstShard)
	err = json.Unmarshal(data, val)
	return val, err
}

func (t *index) SetMeta(id string, data *meta.IndexMeta) error {
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return db.Set(t.key(id, "meta"), val)
}

func (t *index) SetStats(id string, data *meta.IndexStat) error {
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return db.Set(t.key(id, "stats"), val)
}

func (t *index) SetSettings(id string, data *meta.IndexSettings) error {
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return db.Set(t.key(id, "settings"), val)
}

func (t *index) SetMappings(id string, data *meta.Mappings) error {
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return db.Set(t.key(id, "mappings"), val)
}

func (t *index) SetShards(id string, data []*meta.IndexFirstShard) error {
	for _, shard := range data {
		if err := t.SetShard(id, shard); err != nil {
			return err
		}
	}
	return nil
}

func (t *index) SetShard(id string, data *meta.IndexFirstShard) error {
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return db.Set(t.key(id, "shards", data.ID), val)
}

func (t *index) Delete(id string) error {
	for _, key := range []string{"meta", "stats", "settings", "mappings", "shards"} {
		err := db.DeleteWithPrefix(t.key(id, key))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *index) key(keys ...string) string {
	s := new(strings.Builder)
	s.WriteString("/index")
	for _, k := range keys {
		if k != "/" {
			s.WriteString("/")
		}
		s.WriteString(k)
	}
	return s.String()
}
