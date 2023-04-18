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

	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/meta"
)

func TermsQuery(query map[string]interface{}, mappings *meta.Mappings) (bluge.Query, error) {
	if len(query) > 2 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[terms] query doesn't support multiple fields")
	}

	field := ""
	values := []string{}
	valueFloat := []float64{}
	valueInts := []int{}
	valueBools := []bool{}
	boost := -1.0
	for k, v := range query {
		if strings.ToLower(k) == "boost" {
			boost = v.(float64)
			continue
		}

		field = k
		switch v := v.(type) {
		case []string:
			values = v
		case []float64:
			valueFloat = v
		case []int:
			valueInts = v
		case []bool:
			valueBools = v
		case []interface{}:
			for _, vv := range v {
				switch vvv := vv.(type) {
				case string:
					values = append(values, vvv)
				case float64:
					valueFloat = append(valueFloat, vvv)
				case int:
					valueInts = append(valueInts, vvv)
				case bool:
					valueBools = append(valueBools, vvv)
				default:
					return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[term] doesn't support values of type: %T", vv))
				}
			}
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[terms] doesn't support values of type: %T", v))
		}
	}

	subq := bluge.NewBooleanQuery()
	for _, term := range values {
		subqq, err := TermQueryText(field, &meta.TermQuery{Value: term})
		if err != nil {
			return nil, err
		}
		subq.AddShould(subqq)
	}
	for _, term := range valueFloat {
		subqq, err := TermQueryNumeric(field, &meta.TermQuery{Value: term})
		if err != nil {
			return nil, err
		}
		subq.AddShould(subqq)
	}
	for _, term := range valueInts {
		subqq, err := TermQueryNumeric(field, &meta.TermQuery{Value: term})
		if err != nil {
			return nil, err
		}
		subq.AddShould(subqq)
	}
	for _, term := range valueBools {
		subqq, err := TermQueryBool(field, &meta.TermQuery{Value: term})
		if err != nil {
			return nil, err
		}
		subq.AddShould(subqq)
	}
	if boost >= 0 {
		subq.SetBoost(boost)
	}

	return subq, nil
}
