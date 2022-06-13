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

package pebble

import (
	"path"

	"github.com/cockroachdb/pebble"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/metadata/storage"
)

type pebbleStorage struct {
	db *pebble.DB
}

func New(dbpath string) storage.Storager {
	db, err := openpebbleDB(path.Join(config.Global.DataPath, dbpath), false)
	if err != nil {
		log.Fatal().Err(err).Msg("open pebble db for metadata failed")
	}
	return &pebbleStorage{db}
}

func openpebbleDB(dbpath string, readOnly bool) (*pebble.DB, error) {
	opt := &pebble.Options{}
	opt.ReadOnly = readOnly
	opt.Logger = nil
	return pebble.Open(dbpath, opt)
}

func (t *pebbleStorage) List(prefix string, _, _ int) ([][]byte, error) {
	data := make([][]byte, 0)
	iter := t.db.NewIter(t.prefixOption([]byte(prefix)))
	for iter.First(); iter.Valid(); iter.Next() {
		val := iter.Value()
		valCopy := make([]byte, len(val))
		copy(valCopy, val)
		data = append(data, valCopy)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return data, nil
}

func (t *pebbleStorage) Get(key string) ([]byte, error) {
	v, closer, err := t.db.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	data := make([]byte, len(v))
	copy(data, v)
	if err := closer.Close(); err != nil {
		return nil, err
	}
	return data, nil
}

func (t *pebbleStorage) Set(key string, value []byte) error {
	if key == "" {
		return errors.ErrEmptyKey
	}
	return t.db.Set([]byte(key), value, pebble.Sync)
}

func (t *pebbleStorage) Delete(key string) error {
	if key == "" {
		return errors.ErrEmptyKey
	}
	return t.db.Delete([]byte(key), pebble.Sync)
}

func (t *pebbleStorage) Close() error {
	return t.db.Close()
}

func (t *pebbleStorage) prefixOption(prefix []byte) *pebble.IterOptions {
	keyUpperBound := func(b []byte) []byte {
		end := make([]byte, len(b))
		copy(end, b)
		for i := len(end) - 1; i >= 0; i-- {
			end[i] = end[i] + 1
			if end[i] != 0 {
				return end[:i+1]
			}
		}
		return nil // no upper-bound
	}

	return &pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: keyUpperBound(prefix),
	}
}
