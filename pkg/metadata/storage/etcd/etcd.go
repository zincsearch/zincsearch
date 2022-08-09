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
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	client "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/metadata/storage"
)

type etcdStorage struct {
	prefix  string
	timeout time.Duration
	cli     *client.Client
}

// New create an etcd storage instance
// prefix: etcd prefix
// timeout: etcd timeout (default 30s, unit: second)
func New(prefix string, timeout int64) storage.Storager {
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
		prefix:  prefix,
		timeout: time.Second * time.Duration(timeout),
		cli:     cli,
	}
}

func (t *etcdStorage) NewLocker(prefix string) (sync.Locker, error) {
	session, err := concurrency.NewSession(t.cli)
	if err != nil {
		return nil, err
	}

	lock := concurrency.NewLocker(session, t.prefix+prefix)
	return &etcdStorageLocker{
		session: session,
		mutex:   lock,
	}, nil
}

func (t *etcdStorage) List(prefix string, offset, limit int64) ([][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	resp, err := t.cli.Get(ctx, t.prefix+prefix, client.WithPrefix(), client.WithLimit(int64(offset+limit)))
	if err != nil {
		return nil, err
	}
	data := make([][]byte, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		value := make([]byte, len(kv.Value))
		copy(value, kv.Value)
		data = append(data, value)
	}
	if int(offset) >= len(data) {
		return nil, nil
	}
	return data[offset:], nil
}

func (t *etcdStorage) ListEntries(prefix string, offset, limit int64) ([]*storage.StorageEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	prefix = t.prefix + prefix
	resp, err := t.cli.Get(ctx, prefix, client.WithPrefix(), client.WithLimit(int64(offset+limit)))
	if err != nil {
		return nil, err
	}
	data := make([]*storage.StorageEntry, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		entry := &storage.StorageEntry{}
		entry.Key = append(entry.Key, kv.Key[len(prefix):]...)
		entry.Value = append(entry.Value, kv.Value...)
		data = append(data, entry)
	}
	if int(offset) >= len(data) {
		return nil, nil
	}
	return data[offset:], nil
}

func (t *etcdStorage) Get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	_, err := t.cli.Put(ctx, t.prefix+key, string(value))
	return err
}

func (t *etcdStorage) SetWithKeepAlive(key string, value []byte, ttl int64) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}

	lease, _ := t.cli.Lease.Grant(context.Background(), ttl)
	event, err := t.cli.Lease.KeepAlive(context.Background(), lease.ID)
	if err != nil {
		return err
	}
	go func() {
		for range event {
			<-event
		}
		fmt.Println("etcdStorage: keepalive lease expired")
	}()

	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	_, err = t.cli.Put(ctx, t.prefix+key, string(value), client.WithLease(lease.ID))
	return err
}

func (t *etcdStorage) CancelWithKeepAlive(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	resp, err := t.cli.Get(ctx, t.prefix+key)
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		return nil
	}
	_, err = t.cli.Lease.Revoke(context.Background(), client.LeaseID(resp.Kvs[0].Lease))
	return err
}

func (t *etcdStorage) Delete(key string) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	_, err := t.cli.Delete(ctx, t.prefix+key)
	return err
}

func (t *etcdStorage) Watch(key string) <-chan storage.StorageEvent {
	evs := t.cli.Watch(context.Background(), t.prefix+key, client.WithPrefix())
	chs := make(chan storage.StorageEvent, 16)
	bkey := []byte(t.prefix + key)
	go func() {
		for ev := range evs {
			for _, e := range ev.Events {
				var eType int64
				if e.Type == client.EventTypePut {
					eType = meta.StorageEventTypePut
				} else if e.Type == client.EventTypeDelete {
					eType = meta.StorageEventTypeDelete
				}
				chs <- storage.StorageEvent{
					Type:  eType,
					Key:   bytes.TrimPrefix(e.Kv.Key, bkey),
					Value: e.Kv.Value,
				}
			}
		}
	}()
	return chs
}

func (t *etcdStorage) Close() error {
	return t.cli.Close()
}

type etcdStorageLocker struct {
	session *concurrency.Session
	mutex   sync.Locker
}

func (l *etcdStorageLocker) Lock() {
	l.mutex.Lock()
}

func (l *etcdStorageLocker) Unlock() {
	l.mutex.Unlock()
	l.session.Close()
}
