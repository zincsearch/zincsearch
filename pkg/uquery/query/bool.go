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
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func BoolQuery(query map[string]interface{}, mappings *meta.Mappings, analyzers map[string]*analysis.Analyzer) (bluge.Query, error) {
	boolQuery := bluge.NewBooleanQuery()
	var minimumShouldMatch interface{}
	for k, v := range query {
		k := strings.ToLower(k)
		switch k {
		case "should":
			switch v := v.(type) {
			case map[string]interface{}:
				if subq, err := Query(v, mappings, analyzers); err != nil {
					return nil, errors.New(errors.ErrorTypeXContentParseException, "[should] failed to parse field").Cause(err)
				} else {
					boolQuery.AddShould(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := Query(vv.(map[string]interface{}), mappings, analyzers); err != nil {
						return nil, errors.New(errors.ErrorTypeXContentParseException, "[should] failed to parse field").Cause(err)
					} else {
						boolQuery.AddShould(subq)
					}
				}
			default:
				return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[bool] %s doesn't support values of type: %T", k, v))
			}
		case "must":
			switch v := v.(type) {
			case map[string]interface{}:
				if subq, err := Query(v, mappings, analyzers); err != nil {
					return nil, errors.New(errors.ErrorTypeXContentParseException, "[must] failed to parse field").Cause(err)
				} else {
					boolQuery.AddMust(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := Query(vv.(map[string]interface{}), mappings, analyzers); err != nil {
						return nil, errors.New(errors.ErrorTypeXContentParseException, "[must] failed to parse field").Cause(err)
					} else {
						boolQuery.AddMust(subq)
					}
				}
			default:
				return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[bool] %s doesn't support values of type: %T", k, v))
			}
		case "must_not":
			switch v := v.(type) {
			case map[string]interface{}:
				if subq, err := Query(v, mappings, analyzers); err != nil {
					return nil, errors.New(errors.ErrorTypeXContentParseException, "[must_not] failed to parse field").Cause(err)
				} else {
					boolQuery.AddMustNot(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := Query(vv.(map[string]interface{}), mappings, analyzers); err != nil {
						return nil, errors.New(errors.ErrorTypeXContentParseException, "[must_not] failed to parse field").Cause(err)
					} else {
						boolQuery.AddMustNot(subq)
					}
				}
			default:
				return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[bool] %s doesn't support values of type: %T", k, v))
			}
		case "filter":
			filterQuery := bluge.NewBooleanQuery().SetBoost(0)
			switch v := v.(type) {
			case map[string]interface{}:
				if subq, err := Query(v, mappings, analyzers); err != nil {
					return nil, errors.New(errors.ErrorTypeParsingException, "[filter] failed to parse field").Cause(err)
				} else {
					filterQuery.AddMust(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := Query(vv.(map[string]interface{}), mappings, analyzers); err != nil {
						return nil, errors.New(errors.ErrorTypeParsingException, "[filter] failed to parse field").Cause(err)
					} else {
						filterQuery.AddMust(subq)
					}
				}
			default:
				return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[bool] %s doesn't support values of type: %T", k, v))
			}
			boolQuery.AddMust(filterQuery)
		case "minimum_should_match":
			minimumShouldMatch = v
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[bool] unknown field [%s]", k))
		}
	}

	if minimumShouldMatch != nil {
		minValue, err := zutils.CalculateMin(len(boolQuery.Shoulds()), minimumShouldMatch)
		if err != nil {
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[bool] unsupported MinimumShouldMatch value: %v", err))
		}
		boolQuery.SetMinShould(minValue)
	}

	return boolQuery, nil
}
