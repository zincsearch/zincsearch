package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func BoolQuery(query map[string]interface{}, mappings *meta.Mappings, analyzers map[string]*analysis.Analyzer) (bluge.Query, error) {
	boolQuery := bluge.NewBooleanQuery()
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
			boolQuery.SetMinShould(int(v.(float64)))
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[bool] unknown field [%s]", k))
		}
	}

	return boolQuery, nil
}
