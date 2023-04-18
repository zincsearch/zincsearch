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

package timerange

import (
	"strings"

	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/zutils/json"
)

func Query(query interface{}) (int64, int64) {
	if query == nil {
		return 0, 0
	}

	if q, ok := query.(*meta.Query); ok {
		data, err := json.Marshal(q)
		if err != nil {
			return 0, 0
		}
		var newQuery map[string]interface{}
		if err = json.Unmarshal(data, &newQuery); err != nil {
			return 0, 0
		}
		query = newQuery
	}
	q, ok := query.(map[string]interface{})
	if !ok {
		return 0, 0
	}

	for k, t := range q {
		k := strings.ToLower(k)
		v, ok := t.(map[string]interface{})
		if !ok {
			return 0, 0
		}
		switch k {
		case "bool":
			return BoolQuery(v)
		case "range":
			return RangeQuery(v)
		}
	}

	return 0, 0
}
