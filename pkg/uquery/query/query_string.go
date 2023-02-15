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
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/analyzer"
	querystr "github.com/blugelabs/query_string"

	"github.com/zinclabs/zincsearch/pkg/errors"
	"github.com/zinclabs/zincsearch/pkg/meta"
	zincanalysis "github.com/zinclabs/zincsearch/pkg/uquery/analysis"
)

func QueryStringQuery(query map[string]interface{}, mappings *meta.Mappings, analyzers map[string]*analysis.Analyzer) (bluge.Query, error) {
	value := new(meta.QueryStringQuery)
	for k, v := range query {
		k := strings.ToLower(k)
		switch k {
		case "query":
			value.Query = v.(string)
		case "analyzer":
			value.Analyzer = v.(string)
		case "fields":
			if vv, ok := v.([]interface{}); ok {
				for _, vvv := range vv {
					value.Fields = append(value.Fields, vvv.(string))
				}
			}
		case "default_field":
			value.DefaultField = v.(string)
		case "default_operator":
			value.DefaultOperator = v.(string)
		case "boost":
			value.Boost = v.(float64)
		case "analyze_wildcard":
			// noop
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[query_string] unsupported children %s", k))
		}
	}

	options := querystr.DefaultOptions()

	// TODO fields
	// TODO default_field
	// TODO default_operator
	// TODO boost

	zer, _ := zincanalysis.QueryAnalyzer(analyzers, value.Analyzer)
	if zer == nil {
		zer = analyzer.NewStandardAnalyzer()
	}
	options.WithDefaultAnalyzer(zer)

	if len(value.Fields) == 0 {
		for field, prop := range mappings.ListProperty() {
			if prop.Type == "text" || prop.Type == "keyword" {
				value.Fields = append(value.Fields, field)
			}
		}
	}
	for _, field := range value.Fields {
		var zer *analysis.Analyzer
		prop, _ := mappings.GetProperty(field)
		if prop.Type == "keyword" {
			zer = analyzer.NewKeywordAnalyzer()
		} else {
			indexZer, searchZer := zincanalysis.QueryAnalyzerForField(analyzers, mappings, field)
			if zer == nil && searchZer != nil {
				zer = searchZer
			}
			if zer == nil && indexZer != nil {
				zer = indexZer
			}
		}
		if zer != nil {
			options.WithAnalyzerForField(field, zer)
		}
	}

	return querystr.ParseQueryString(value.Query, options)
}
