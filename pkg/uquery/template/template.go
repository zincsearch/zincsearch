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

package template

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zinclabs/zincsearch/pkg/errors"
	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/pkg/uquery/index"
)

func Request(data map[string]interface{}) (*meta.IndexTemplate, error) {
	if data == nil {
		return nil, nil
	}

	if data["template"] == nil {
		return nil, errors.New(errors.ErrorTypeXContentParseException, "[template] template should be defined")
	}

	template := new(meta.IndexTemplate)
	for k, v := range data {
		k = strings.ToLower(k)
		switch k {
		case "name", "data_stream":
			// ignore
		case "index_patterns":
			patterns, ok := v.([]interface{})
			if !ok {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[template] index_patterns value should be an array of string")
			}
			for _, pattern := range patterns {
				template.IndexPatterns = append(template.IndexPatterns, pattern.(string))
			}
		case "priority":
			switch v := v.(type) {
			case string:
				template.Priority, _ = strconv.Atoi(v)
			case float64:
				template.Priority = int(v)
			case int:
				template.Priority = v
			default:
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[template] priority value should be a numberic")
			}
		case "template":
			switch v := v.(type) {
			case string:
				// compatible {"priority":150,"template":"filebeat-7.16.3-*"}
				template.IndexPatterns = append(template.IndexPatterns, v)
			case map[string]interface{}:
				tmpIndex, err := index.Request(v)
				if err != nil {
					return nil, err
				}
				template.Template.Settings = tmpIndex.Settings
				template.Template.Mappings = tmpIndex.Mappings
			default:
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[template] template value should be an object")
			}
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[template] unknown option [%s]", k))
		}
	}

	if len(template.IndexPatterns) == 0 {
		return nil, errors.New(errors.ErrorTypeXContentParseException, "[template] index_patterns should be defined")
	}

	return template, nil
}
