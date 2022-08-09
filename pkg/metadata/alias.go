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
	"github.com/goccy/go-json"
)

type alias struct{}

var Alias = new(alias)

const aliasKey = "/aliases/alias"

func (t *alias) Set(data map[string][]string) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return db.Set(aliasKey, buf)
}

func (t *alias) Get() (map[string][]string, error) {
	data, err := db.Get(aliasKey)
	if err != nil {
		return nil, err
	}

	als := map[string][]string{}
	err = json.Unmarshal(data, als)
	return als, err
}
