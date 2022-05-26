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

type kv struct{}

var KV = new(kv)

func (t *kv) List(offset, limit int) ([][]byte, error) {
	return db.List(t.key(""), offset, limit)
}

func (t *kv) Get(key string) ([]byte, error) {
	return db.Get(t.key(key))
}

func (t *kv) Set(key string, val []byte) error {
	return db.Set(t.key(key), val)
}

func (t *kv) Delete(key string) error {
	return db.Delete(t.key(key))
}

func (t *kv) key(key string) string {
	return "/kv/" + key
}
