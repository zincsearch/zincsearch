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

	"github.com/zinclabs/zinc/pkg/config"
	"github.com/zinclabs/zinc/pkg/metadata/storage"
	"github.com/zinclabs/zinc/pkg/metadata/storage/badger"
	"github.com/zinclabs/zinc/pkg/metadata/storage/bolt"
	"github.com/zinclabs/zinc/pkg/metadata/storage/etcd"
)

var DefaultTimeout int64 = 30 // 30s

var db storage.Storager

func init() {
	switch strings.ToLower(config.Global.MetadataStorage) {
	case "badger":
		db = badger.New("_metadata.db")
	case "bolt":
		db = bolt.New("_metadata.bolt")
	case "etcd":
		prefix := strings.TrimRight(config.Global.Etcd.Prefix, "/")
		db = etcd.New(prefix, DefaultTimeout)
	default:
		db = bolt.New("_metadata.bolt")
	}
}

func Close() error {
	return db.Close()
}
