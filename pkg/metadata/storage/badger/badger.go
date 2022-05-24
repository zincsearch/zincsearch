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

package badger

import (
	"path"
	"runtime"

	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/metadata/storage"
)

type badgerStorage struct {
	db *badger.DB
}

func New(dbpath string) storage.Storager {
	db, err := openBadgerDB(path.Join(config.Global.DataPath, dbpath), false)
	if err != nil {
		log.Fatal().Err(err).Msg("open badger db for metadata failed")
	}
	return &badgerStorage{db}
}

func openBadgerDB(dbpath string, readOnly bool) (*badger.DB, error) {
	opt := badger.DefaultOptions(dbpath)
	opt.NumGoroutines = runtime.NumGoroutine() * 8
	opt.MemTableSize = 32 << 20
	opt.Compression = options.ZSTD
	opt.ZSTDCompressionLevel = 3
	opt.BlockSize = 1024 * 128
	opt.MetricsEnabled = false
	// opt.Logger = nil
	opt.ReadOnly = readOnly
	return badger.Open(opt)
}

func (t *badgerStorage) List(prefix string, offset, limit int) ([][]byte, error) {
	return nil, nil
}

func (t *badgerStorage) Get(key string) ([]byte, error) {
	return nil, nil
}

func (t *badgerStorage) Set(key string, value []byte) error {
	return nil
}

func (t *badgerStorage) Delete(key string) error {
	return nil
}

func (t *badgerStorage) Close() error {
	return t.db.Close()
}
