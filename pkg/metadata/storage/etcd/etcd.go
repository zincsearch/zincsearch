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

package etcd

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	client "go.etcd.io/etcd/client/v3"

	"github.com/zinclabs/zincsearch/pkg/config"
	"github.com/zinclabs/zincsearch/pkg/errors"
	"github.com/zinclabs/zincsearch/pkg/metadata/storage"
)

var timeout = 30 * time.Second

type etcdStorage struct {
	prefix string
	cli    *client.Client
}

func New(dbpath string) storage.Storager {
	cli, err := client.New(client.Config{
		Endpoints:   config.Global.Etcd.Endpoints,
		DialTimeout: 5 * time.Second,
		Username:    config.Global.Etcd.Username,
		Password:    config.Global.Etcd.Password,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("open etcd for metadata failed")
	}
	return &etcdStorage{
		prefix: dbpath,
		cli:    cli,
	}
}

func (t *etcdStorage) List(prefix string, _, _ int) ([][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := t.cli.Get(ctx, t.prefix+prefix, client.WithPrefix())
	if err != nil {
		return nil, err
	}
	data := make([][]byte, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		data = append(data, kv.Value)
	}
	return data, nil
}

func (t *etcdStorage) Get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := t.cli.Get(ctx, t.prefix+key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, errors.ErrKeyNotFound
	}
	return resp.Kvs[0].Value, nil
}

func (t *etcdStorage) Set(key string, value []byte) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := t.cli.Put(ctx, t.prefix+key, string(value))
	return err
}

func (t *etcdStorage) Delete(key string) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := t.cli.Delete(ctx, t.prefix+key)
	return err
}

func (t *etcdStorage) Close() error {
	return t.cli.Close()
}
