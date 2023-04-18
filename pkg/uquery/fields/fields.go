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

package fields

import (
	"strings"
	"time"

	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/zutils/json"
)

func Request(v []interface{}) ([]*meta.Field, error) {
	if v == nil {
		return nil, nil
	}

	fields := make([]*meta.Field, 0, len(v))
	for _, v := range v {
		switch v := v.(type) {
		case string:
			fields = append(fields, &meta.Field{Field: v})
		case map[string]interface{}:
			f := new(meta.Field)
			for k, v := range v {
				k = strings.ToLower(k)
				switch k {
				case "field":
					f.Field = v.(string)
				case "format":
					f.Format = v.(string)
				default:
					// return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[fields] unknown field [%s]", k))
				}
			}
			fields = append(fields, f)
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, "[fields] value should be string or object")

		}
	}

	return fields, nil
}

func Response(fields []*meta.Field, data []byte, mappings *meta.Mappings) map[string]interface{} {
	// return empty
	if len(fields) == 0 {
		return nil
	}

	ret := make(map[string]interface{})
	err := json.Unmarshal(data, &ret)
	if err != nil {
		return nil
	}

	var field string
	wildcard := false
	results := make(map[string]interface{})
	for _, v := range fields {
		wildcard = false
		field = v.Field
		if strings.HasSuffix(field, "*") {
			wildcard = true
		}
		if rv, ok := ret[field]; ok {
			prop, _ := mappings.GetProperty(field)
			if (prop.Type == "date" || prop.Type == "time") && v.Format != "" {
				if t, err := time.Parse(prop.Format, rv.(string)); err == nil {
					results[field] = []interface{}{t.Format(v.Format)}
				} else {
					results[field] = []interface{}{rv}
				}
			} else {
				results[field] = []interface{}{rv}
			}
		} else if wildcard {
			for rk, rv := range ret {
				if strings.HasPrefix(rk, field[:len(field)-1]) {
					prop, _ := mappings.GetProperty(rk)
					if (prop.Type == "date" || prop.Type == "time") && v.Format != "" {
						if t, err := time.Parse(prop.Format, rv.(string)); err == nil {
							results[rk] = []interface{}{t.Format(v.Format)}
						} else {
							results[rk] = []interface{}{rv}
						}
					} else {
						results[rk] = []interface{}{rv}
					}
				}
			}
		}
	}

	return results
}
