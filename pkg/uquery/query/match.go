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

	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/meta"
	zincanalysis "github.com/zincsearch/zincsearch/pkg/uquery/analysis"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func MatchQuery(query map[string]interface{}, mappings *meta.Mappings, analyzers map[string]*analysis.Analyzer) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[match] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.MatchQuery)
	value.Boost = -1.0
	for k, v := range query {
		field = k
		switch v := v.(type) {
		case string:
			value.Query = v
		case map[string]interface{}:
			for k, v := range v {
				k := strings.ToLower(k)
				switch k {
				case "query":
					value.Query, _ = zutils.ToString(v)
				case "analyzer":
					value.Analyzer, _ = zutils.ToString(v)
				case "operator":
					value.Operator, _ = zutils.ToString(v)
				case "fuzziness":
					value.Fuzziness = v
				case "prefix_length":
					value.PrefixLength, _ = zutils.ToFloat64(v)
				case "boost":
					value.Boost, _ = zutils.ToFloat64(v)
				default:
					// return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[match] unknown field [%s]", k))
				}
			}
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[match] %s doesn't support values of type: %T", k, v))
		}
	}

	var err error
	var zer *analysis.Analyzer
	if value.Analyzer != "" {
		zer, err = zincanalysis.QueryAnalyzer(analyzers, value.Analyzer)
		if err != nil {
			return nil, err
		}
	} else {
		indexZer, searchZer := zincanalysis.QueryAnalyzerForField(analyzers, mappings, field)
		if zer == nil && searchZer != nil {
			zer = searchZer
		}
		if zer == nil && indexZer != nil {
			zer = indexZer
		}
	}

	subq := bluge.NewMatchQuery(value.Query).SetField(field)
	if zer != nil {
		subq.SetAnalyzer(zer)
	}
	if value.Operator != "" {
		op := strings.ToUpper(value.Operator)
		switch op {
		case "OR":
			subq.SetOperator(bluge.MatchQueryOperatorOr)
		case "AND":
			subq.SetOperator(bluge.MatchQueryOperatorAnd)
		default:
			return nil, errors.New(errors.ErrorTypeIllegalArgumentException, fmt.Sprintf("[match] unknown operator %s", op))
		}
	}
	if value.Fuzziness != nil {
		if value.Fuzziness != nil {
			v := ParseFuzziness(value.Fuzziness, value.Query, zer)
			if v > 0 {
				subq.SetFuzziness(v)
			}
		}
	}
	if value.PrefixLength > 0 {
		subq.SetPrefix(int(value.PrefixLength))
	}
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}
