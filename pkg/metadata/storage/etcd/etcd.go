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
	"time"

	"github.com/rs/zerolog/log"
	client "go.etcd.io/etcd/client/v3"

	"github.com/zinclabs/zinc/pkg/metadata/storage"
)

type etcdStorage struct {
	prefix string
	cli    *client.Client
}

func New(dbpath string) storage.Storager {
	cli, err := client.New(client.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("open etcd for metadata failed")
	}
	return &etcdStorage{
		prefix: dbpath,
		cli:    cli,
	}
}

func (t *etcdStorage) List(prefix string, offset, limit int) ([][]byte, error) {
	return nil, nil
}

func (t *etcdStorage) Get(key string) ([]byte, error) {
	return nil, nil
}

func (t *etcdStorage) Set(key string, value []byte) error {
	return nil
}

func (t *etcdStorage) Delete(key string) error {
	return nil
}

func (t *etcdStorage) Close() error {
	return t.cli.Close()
}
