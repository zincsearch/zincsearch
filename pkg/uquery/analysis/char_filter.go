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

package analysis

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/char"

	"github.com/zinclabs/zincsearch/pkg/errors"
	zincchar "github.com/zinclabs/zincsearch/pkg/uquery/analysis/char"
	"github.com/zinclabs/zincsearch/pkg/zutils"
)

func RequestCharFilter(data map[string]interface{}) (map[string]analysis.CharFilter, error) {
	if data == nil {
		return nil, nil
	}

	filters := make(map[string]analysis.CharFilter)
	for name, options := range data {
		typ, err := zutils.GetStringFromMap(options, "type")
		if err != nil {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[char_filter] %s option [%s] should be exists", name, "type"))
		}
		filter, err := RequestCharFilterSingle(typ, options)
		if err != nil {
			return nil, err
		}
		filters[name] = filter
	}

	return filters, nil
}

func RequestCharFilterSlice(data []interface{}) ([]analysis.CharFilter, error) {
	if data == nil {
		return nil, nil
	}

	filters := make([]analysis.CharFilter, 0, len(data))
	for _, options := range data {
		var err error
		var filter analysis.CharFilter
		switch v := options.(type) {
		case string:
			filter, err = RequestCharFilterSingle(v, nil)
		case map[string]interface{}:
			var typ string
			typ, err = zutils.GetStringFromMap(options, "type")
			if err != nil {
				return nil, errors.New(errors.ErrorTypeParsingException, "[char_filter] option [type] should be exists")
			}
			filter, err = RequestCharFilterSingle(typ, options)
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, "[char_filter] option should be string or object")
		}
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}

	return filters, nil
}

func RequestCharFilterSingle(name string, options interface{}) (analysis.CharFilter, error) {
	name = strings.ToLower(name)
	switch name {
	case "ascii_folding", "asciifolding":
		return char.NewASCIIFoldingFilter(), nil
	case "html", "html_strip":
		return char.NewHTMLCharFilter(), nil
	case "zero_width_non_joiner":
		return char.NewZeroWidthNonJoinerCharFilter(), nil
	case "regexp", "pattern", "pattern_replace":
		return zincchar.NewRegexpCharFilter(options)
	case "mapping":
		return zincchar.NewMappingCharFilter(options)
	default:
		return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[char_filter] unknown character filter [%s]", name))
	}
}
