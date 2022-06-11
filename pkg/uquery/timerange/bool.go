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

import "strings"

func BoolQuery(query map[string]interface{}) (int64, int64) {
	for k, v := range query {
		k := strings.ToLower(k)
		switch k {
		case "should":
			switch v := v.(type) {
			case map[string]interface{}:
				return Query(v)
			case []interface{}:
				for _, vv := range v {
					min, max := Query(vv.(map[string]interface{}))
					if min > 0 || max > 0 {
						return min, max
					}
				}
			}
		case "must":
			switch v := v.(type) {
			case map[string]interface{}:
				return Query(v)
			case []interface{}:
				for _, vv := range v {
					min, max := Query(vv.(map[string]interface{}))
					if min > 0 || max > 0 {
						return min, max
					}
				}
			}
		case "must_not":
		case "filter":
			switch v := v.(type) {
			case map[string]interface{}:
				return Query(v)
			case []interface{}:
				for _, vv := range v {
					min, max := Query(vv.(map[string]interface{}))
					if min > 0 || max > 0 {
						return min, max
					}
				}
			}
		}
	}

	return 0, 0
}
