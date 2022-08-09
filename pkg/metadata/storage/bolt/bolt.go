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

package bolt

import (
	"bytes"
	"os"
	"path"
	"sync"

	"github.com/rs/zerolog/log"
	"go.etcd.io/bbolt"

	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/metadata/storage"
)

type boltStorage struct {
	db   *bbolt.DB
	lock sync.Map
}

func New(dbpath string) storage.Storager {
	db, err := openbboltDB(path.Join(config.Global.DataPath, dbpath), false)
	if err != nil {
		log.Fatal().Err(err).Msg("open bbolt db for metadata failed")
	}
	return &boltStorage{db: db}
}

func openbboltDB(dbpath string, readOnly bool) (*bbolt.DB, error) {
	opt := &bbolt.Options{
		Timeout:      0,
		NoGrowSync:   false,
		FreelistType: bbolt.FreelistArrayType,
	}
	if err := os.MkdirAll(path.Dir(dbpath), 0755); err != nil {
		return nil, err
	}
	return bbolt.Open(dbpath, 0666, opt)
}

func (t *boltStorage) NewLocker(prefix string) (sync.Locker, error) {
	if lock, ok := t.lock.Load(prefix); ok {
		return lock.(*boltStorageLocker), nil
	}
	lock := &boltStorageLocker{key: prefix, db: t, mutex: &sync.Mutex{}}
	t.lock.Store(prefix, lock)
	return lock, nil
}

func (t *boltStorage) List(prefix string, offset, limit int64) ([][]byte, error) {
	data := make([][]byte, 0)
	bucket, _ := t.splitBucketAndKey(prefix)
	err := t.db.View(func(txn *bbolt.Tx) error {
		b := txn.Bucket(bucket)
		if b == nil {
			return nil
		}
		i := int64(0)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			i++
			if i <= offset {
				continue
			}
			if limit > 0 && i > offset+limit {
				break
			}
			valCopy := make([]byte, len(v))
			copy(valCopy, v)
			data = append(data, valCopy)
		}
		return nil
	})
	return data, err
}

func (t *boltStorage) ListEntries(prefix string, offset, limit int64) ([]*storage.StorageEntry, error) {
	prefixByte := []byte(prefix)
	data := make([]*storage.StorageEntry, 0)
	bucket, _ := t.splitBucketAndKey(prefix)
	err := t.db.View(func(txn *bbolt.Tx) error {
		b := txn.Bucket(bucket)
		if b == nil {
			return nil
		}
		i := int64(0)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			i++
			if i <= offset {
				continue
			}
			if limit > 0 && i > offset+limit {
				break
			}
			entry := &storage.StorageEntry{}
			entry.Key = bytes.TrimPrefix(k, prefixByte)
			entry.Value = make([]byte, len(v))
			copy(entry.Value, v)
			data = append(data, entry)
		}
		return nil
	})
	return data, err
}

func (t *boltStorage) Get(key string) ([]byte, error) {
	var data []byte
	bucket, name := t.splitBucketAndKey(key)
	err := t.db.View(func(txn *bbolt.Tx) error {
		b := txn.Bucket(bucket)
		if b == nil {
			return errors.ErrKeyNotFound
		}
		v := b.Get(name)
		if v == nil {
			return errors.ErrKeyNotFound
		}
		data = make([]byte, len(v))
		copy(data, v)
		return nil
	})
	return data, err
}

func (t *boltStorage) Set(key string, value []byte) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	bucket, name := t.splitBucketAndKey(key)
	return t.db.Update(func(txn *bbolt.Tx) error {
		b, err := txn.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return b.Put(name, value)
	})
}

func (t *boltStorage) SetWithKeepAlive(key string, value []byte, _ int64) error {
	return t.Set(key, value)
}

func (t *boltStorage) CancelWithKeepAlive(key string) error {
	return t.Delete(key)
}

func (t *boltStorage) Delete(key string) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	bucket, name := t.splitBucketAndKey(key)
	return t.db.Update(func(Tx *bbolt.Tx) error {
		b := Tx.Bucket(bucket)
		if b == nil {
			return nil
		}
		return b.Delete(name)
	})
}

func (t *boltStorage) Watch(key string) <-chan storage.StorageEvent {
	return make(chan storage.StorageEvent, 1)
}

func (t *boltStorage) Close() error {
	return t.db.Close()
}

func (t *boltStorage) splitBucketAndKey(key string) ([]byte, []byte) {
	if key == "" {
		return nil, nil
	}
	p := bytes.LastIndex([]byte(key), []byte("/"))
	return []byte(key[:p]), []byte(key[p+1:])
}

type boltStorageLocker struct {
	key   string
	db    *boltStorage
	mutex sync.Locker
}

func (l *boltStorageLocker) Lock() {
	l.mutex.Lock()
}

func (l *boltStorageLocker) Unlock() {
	l.mutex.Unlock()
	l.db.lock.Delete(l.key)
}
