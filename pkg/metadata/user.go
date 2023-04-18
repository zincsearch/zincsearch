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
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/zutils/json"
)

type user struct{}

var User = new(user)

func (t *user) List(offset, limit int) ([]*meta.User, error) {
	data, err := db.List(t.key(""), offset, limit)
	if err != nil {
		return nil, err
	}
	users := make([]*meta.User, 0, len(data))
	for _, d := range data {
		u := new(meta.User)
		err = json.Unmarshal(d, u)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (t *user) Get(id string) (*meta.User, error) {
	data, err := db.Get(t.key(id))
	if err != nil {
		return nil, err
	}
	u := new(meta.User)
	err = json.Unmarshal(data, u)
	return u, err
}

func (t *user) Set(id string, val meta.User) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return db.Set(t.key(id), data)
}

func (t *user) Delete(id string) error {
	return db.Delete(t.key(id))
}

func (t *user) key(id string) string {
	return "/user/" + id
}
