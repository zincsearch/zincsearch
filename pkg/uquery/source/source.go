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

package source

import (
	"strings"

	"github.com/zinclabs/zincsearch/pkg/errors"
	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/pkg/zutils/json"
)

func Request(v interface{}) (*meta.Source, error) {
	source := &meta.Source{Enable: true}
	if v == nil {
		return source, nil
	}
	if v, ok := v.(*meta.Source); ok {
		return v, nil
	}

	switch v := v.(type) {
	case bool:
		source.Enable = v
	case []interface{}:
		source.Fields = make([]string, 0, len(v))
		for _, field := range v {
			if v, ok := field.(string); ok {
				source.Fields = append(source.Fields, v)
			} else {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[_source] value should be boolean or []string")
			}
		}
	default:
		return nil, errors.New(errors.ErrorTypeXContentParseException, "[_source] value should be boolean or []string")
	}

	return source, nil
}

func Response(source *meta.Source, data []byte) map[string]interface{} {
	// return empty
	if !source.Enable {
		return nil
	}

	ret := make(map[string]interface{})
	err := json.Unmarshal(data, &ret)
	if err != nil {
		return nil
	}

	// return all fields
	if len(source.Fields) == 0 {
		return ret
	}

	wildcard := false
	rets := make(map[string]interface{})
	for _, field := range source.Fields {
		wildcard = false
		if strings.HasSuffix(field, "*") {
			wildcard = true
		}
		if _, ok := ret[field]; ok {
			rets[field] = ret[field]
		} else if wildcard {
			for k, v := range ret {
				if strings.HasPrefix(k, field[:len(field)-1]) {
					rets[k] = v
				}
			}
		}
	}

	return rets
}
