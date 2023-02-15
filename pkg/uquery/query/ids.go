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

package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"

	"github.com/zinclabs/zincsearch/pkg/errors"
	"github.com/zinclabs/zincsearch/pkg/meta"
)

func IdsQuery(query map[string]interface{}, mappings *meta.Mappings) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[ids] query doesn't support multiple fields")
	}

	value := new(meta.IdsQuery)
	for k, v := range query {
		switch v := v.(type) {
		case []string:
			value.Values = v
		case []interface{}:
			value.Values = make([]string, len(v))
			for i, v := range v {
				value.Values[i] = v.(string)
			}
		case map[string]interface{}:
			for k, v := range v {
				k := strings.ToLower(k)
				switch k {
				case "value":
					switch v := v.(type) {
					case []interface{}:
						value.Values = make([]string, len(v))
						for i, v := range v {
							value.Values[i] = v.(string)
						}
					default:
						return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[ids] %s doesn't support values of type: %T", k, v))
					}
				default:
					// return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[ids] unknown field [%s]", k))
				}
			}
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[ids] %s doesn't support values of type: %T", k, v))
		}
	}

	return TermsQuery(map[string]interface{}{
		"_id": value.Values,
	}, mappings)
}
