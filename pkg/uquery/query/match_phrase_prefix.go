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

	"github.com/zinclabs/zincsearch/pkg/errors"
	"github.com/zinclabs/zincsearch/pkg/meta"
	zincanalysis "github.com/zinclabs/zincsearch/pkg/uquery/analysis"
)

func MatchPhrasePrefixQuery(query map[string]interface{}, mappings *meta.Mappings, analyzers map[string]*analysis.Analyzer) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[match_phrase_prefix] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.MatchPhrasePrefixQuery)
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
					value.Query = v.(string)
				case "analyzer":
					value.Analyzer = v.(string)
				case "boost":
					value.Boost = v.(float64)
				default:
					// return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[match_phrase_prefix] unknown field [%s]", k))
				}
			}
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[match_phrase_prefix] %s doesn't support values of type: %T", k, v))
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
	if zer == nil {
		zer = analyzer.NewStandardAnalyzer()
	}

	tokens := zer.Analyze([]byte(value.Query))
	subq := bluge.NewBooleanQuery()
	if len(tokens) > 0 {
		subq.AddMust(bluge.NewPrefixQuery(string(tokens[len(tokens)-1].Term)).SetField(field))
	}
	if len(tokens) > 1 {
		phrase := strings.Builder{}
		for i := 0; i < len(tokens)-1; i++ {
			phrase.WriteString(string(tokens[i].Term))
			phrase.WriteString(" ")
		}
		subq.AddMust(bluge.NewMatchPhraseQuery(strings.TrimSpace(phrase.String())).SetField(field).SetAnalyzer(zer))
	}
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}
