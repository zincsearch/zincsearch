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
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/zutils/json"
)

type role struct{}

var Role = new(role)

func (t *role) List(offset, limit int) ([]*meta.Role, error) {
	data, err := db.List(t.key(""), offset, limit)
	if err != nil {
		return nil, err
	}
	roles := make([]*meta.Role, 0, len(data))
	for _, d := range data {
		u := new(meta.Role)
		err = json.Unmarshal(d, u)
		if err != nil {
			return nil, err
		}
		roles = append(roles, u)
	}
	return roles, nil
}

func (t *role) Get(id string) (*meta.Role, error) {
	data, err := db.Get(t.key(id))
	if err != nil {
		return nil, err
	}
	r := new(meta.Role)
	err = json.Unmarshal(data, r)
	return r, err
}

func (t *role) Set(id string, val meta.Role) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return db.Set(t.key(id), data)
}

func (t *role) Delete(id string) error {
	return db.Delete(t.key(id))
}

func (t *role) key(id string) string {
	return "/role/" + id
}
