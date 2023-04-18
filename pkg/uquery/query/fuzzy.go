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

	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/meta"
	zincanalysis "github.com/zincsearch/zincsearch/pkg/uquery/analysis"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func FuzzyQuery(query map[string]interface{}, mappings *meta.Mappings, analyzers map[string]*analysis.Analyzer) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[fuzzy] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.FuzzyQuery)
	value.Boost = -1.0
	for k, v := range query {
		field = k
		switch v := v.(type) {
		case string:
			value.Value = v
		case map[string]interface{}:
			for k, v := range v {
				k := strings.ToLower(k)
				switch k {
				case "value":
					value.Value, _ = zutils.ToString(v)
				case "fuzziness":
					value.Fuzziness = v
				case "prefix_length":
					value.PrefixLength, _ = zutils.ToFloat64(v)
				case "boost":
					value.Boost, _ = zutils.ToFloat64(v)
				default:
					// return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[fuzzy] unknown field [%s]", k))
				}
			}
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[fuzzy] %s doesn't support values of type: %T", k, v))
		}
	}

	var zer *analysis.Analyzer
	indexZer, searchZer := zincanalysis.QueryAnalyzerForField(analyzers, mappings, field)
	if zer == nil && searchZer != nil {
		zer = searchZer
	}
	if zer == nil && indexZer != nil {
		zer = indexZer
	}

	subq := bluge.NewFuzzyQuery(value.Value).SetField(field)
	if value.Fuzziness != nil {
		v := ParseFuzziness(value.Fuzziness, value.Value, zer)
		if v > 0 {
			subq.SetFuzziness(v)
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

func ParseFuzziness(fuzziness interface{}, query string, zer *analysis.Analyzer) int {
	val, _ := zutils.ToString(fuzziness)
	val = strings.ToUpper(val)
	if !strings.HasPrefix(val, "AUTO") {
		v, _ := zutils.ToInt(val)
		return v
	}

	if zer == nil {
		zer = analyzer.NewStandardAnalyzer()
	}
	tokens := zer.Analyze([]byte(query))
	n := 0
	for _, token := range tokens {
		if n < len(token.Term) {
			n = len(token.Term)
		}
	}

	n1 := 3
	n2 := 6
	if strings.Contains(val, ":") && strings.Contains(val, ",") {
		val := strings.TrimPrefix(val, "AUTO:")
		vals := strings.Split(val, ",")
		if len(vals) == 2 {
			n1, _ = zutils.ToInt(vals[0])
			n2, _ = zutils.ToInt(vals[1])
			if n1 < 2 || n1 >= n2 {
				return 0
			}
		}
	}

	v := 0
	if n >= n2 {
		v = 2
	} else if n >= n1 {
		v = 1
	}
	return v
}
