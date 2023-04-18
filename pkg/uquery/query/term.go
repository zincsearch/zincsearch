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
	"strconv"
	"strings"

	"github.com/blugelabs/bluge"

	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func TermQuery(query map[string]interface{}, mappings *meta.Mappings) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[term] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.TermQuery)
	value.Boost = -1.0
	for k, v := range query {
		field = k
		switch v := v.(type) {
		case string:
			value.Value = v
		case float64:
			value.Value = v
		case bool:
			value.Value = v
		case map[string]interface{}:
			for k, v := range v {
				k := strings.ToLower(k)
				switch k {
				case "value":
					value.Value = v
				case "case_insensitive":
					value.CaseInsensitive = v.(bool)
				case "boost":
					value.Boost = v.(float64)
				default:
					// return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[term] unknown field [%s]", k))
				}
			}
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[term] doesn't support values of type: %T", v))
		}
	}

	// TODO: case_insensitive support

	prop, _ := mappings.GetProperty(field)
	switch prop.Type {
	case "numeric":
		return TermQueryNumeric(field, value)
	case "bool":
		return TermQueryBool(field, value)
	default:
		return TermQueryText(field, value)
	}
}

func TermQueryNumeric(field string, value *meta.TermQuery) (bluge.Query, error) {
	val, err := zutils.ToFloat64(value.Value)
	if err != nil {
		return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[term] convert value to numeric error: %s", err))
	}
	subq := bluge.NewNumericRangeInclusiveQuery(val, val, true, true).SetField(field)
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}
	return subq, nil
}

func TermQueryBool(field string, value *meta.TermQuery) (bluge.Query, error) {
	val, err := zutils.ToBool(value.Value)
	if err != nil {
		return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[term] convert value to boolean error: %s", err))
	}
	subq := bluge.NewTermQuery(strconv.FormatBool(val)).SetField(field)
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}
	return subq, nil
}

func TermQueryText(field string, value *meta.TermQuery) (bluge.Query, error) {
	val, err := zutils.ToString(value.Value)
	if err != nil {
		return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[term] convert value to string error: %s", err))
	}
	subq := bluge.NewTermQuery(val).SetField(field)
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}
	return subq, nil
}
