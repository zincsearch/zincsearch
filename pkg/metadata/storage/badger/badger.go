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
	"bytes"
	"path"
	"sync"

	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/metadata/storage"
)

type badgerStorage struct {
	db   *badger.DB
	lock sync.Map
}

func New(dbpath string) storage.Storager {
	db, err := openBadgerDB(path.Join(config.Global.DataPath, dbpath), false)
	if err != nil {
		log.Fatal().Err(err).Msg("open badger db for metadata failed")
	}
	return &badgerStorage{db: db}
}

func openBadgerDB(dbpath string, readOnly bool) (*badger.DB, error) {
	opt := badger.DefaultOptions(dbpath)
	opt.NumGoroutines = 4
	opt.MemTableSize = 1 << 24 // 16MB
	opt.Compression = options.ZSTD
	opt.ZSTDCompressionLevel = 3
	opt.BlockSize = 4096           // 4KB
	opt.ValueLogFileSize = 1 << 25 // 32MB
	opt.MetricsEnabled = false
	opt.Logger = nil
	opt.ReadOnly = readOnly
	return badger.Open(opt)
}

func (t *badgerStorage) NewLocker(prefix string) (sync.Locker, error) {
	if lock, ok := t.lock.Load(prefix); ok {
		return lock.(*badgerStorageLocker), nil
	}
	lock := &badgerStorageLocker{key: prefix, db: t, mutex: &sync.Mutex{}}
	t.lock.Store(prefix, lock)
	return lock, nil
}

func (t *badgerStorage) List(prefix string, offset, limit int64) ([][]byte, error) {
	data := make([][]byte, 0)
	pre := []byte(prefix)
	err := t.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()
		i := int64(0)
		for it.Seek(pre); it.ValidForPrefix(pre); it.Next() {
			i++
			if i <= offset {
				continue
			}
			if limit > 0 && i > offset+limit {
				break
			}
			item := it.Item()
			buf, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			data = append(data, buf)
		}
		return nil
	})
	return data, err
}

func (t *badgerStorage) ListEntries(prefix string, offset, limit int64) ([]*storage.StorageEntry, error) {
	data := make([]*storage.StorageEntry, 0)
	pre := []byte(prefix)
	err := t.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()
		var err error
		i := int64(0)
		for it.Seek(pre); it.ValidForPrefix(pre); it.Next() {
			i++
			if i <= offset {
				continue
			}
			if limit > 0 && i > offset+limit {
				break
			}
			item := it.Item()
			entry := &storage.StorageEntry{}
			entry.Key = bytes.TrimPrefix(item.Key(), pre)
			entry.Value, err = item.ValueCopy(nil)
			if err != nil {
				return err
			}
			data = append(data, entry)
		}
		return nil
	})
	return data, err
}

func (t *badgerStorage) Get(key string) ([]byte, error) {
	var data []byte
	err := t.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		data, err = item.ValueCopy(nil)
		return err
	})
	if err == badger.ErrKeyNotFound {
		return nil, errors.ErrKeyNotFound
	}
	return data, err
}

func (t *badgerStorage) Set(key string, value []byte) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	return t.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

func (t *badgerStorage) SetWithKeepAlive(key string, value []byte, _ int64) error {
	return t.Set(key, value)
}

func (t *badgerStorage) CancelWithKeepAlive(key string) error {
	return t.Delete(key)
}

func (t *badgerStorage) Delete(key string) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	return t.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (t *badgerStorage) Watch(key string) <-chan storage.StorageEvent {
	return make(chan storage.StorageEvent, 1)
}

func (t *badgerStorage) Close() error {
	return t.db.Close()
}

type badgerStorageLocker struct {
	key   string
	db    *badgerStorage
	mutex sync.Locker
}

func (l *badgerStorageLocker) Lock() {
	l.mutex.Lock()
}

func (l *badgerStorageLocker) Unlock() {
	l.mutex.Unlock()
	l.db.lock.Delete(l.key)
}
