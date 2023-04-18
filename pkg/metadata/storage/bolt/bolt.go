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

	"github.com/rs/zerolog/log"
	"go.etcd.io/bbolt"

	"github.com/zincsearch/zincsearch/pkg/config"
	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/metadata/storage"
)

type boltStorage struct {
	db *bbolt.DB
}

func New(dbpath string) storage.Storager {
	db, err := openbboltDB(path.Join(config.Global.DataPath, dbpath), false)
	if err != nil {
		log.Fatal().Err(err).Msg("open bbolt db for metadata failed")
	}
	return &boltStorage{db}
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

func (t *boltStorage) List(prefix string, _, _ int) ([][]byte, error) {
	data := make([][]byte, 0)
	bucket, _ := t.splitBucketAndKey(prefix)
	err := t.db.View(func(txn *bbolt.Tx) error {
		b := txn.Bucket(bucket)
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			valCopy := make([]byte, len(v))
			copy(valCopy, v)
			data = append(data, valCopy)
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
