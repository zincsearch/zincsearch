package parser

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func BoolQuery(query map[string]interface{}) (bluge.Query, error) {
	boolQuery := bluge.NewBooleanQuery()
	for k, v := range query {
		k := strings.ToLower(k)
		switch k {
		case "should":
			switch v := v.(type) {
			case map[string]interface{}:
				if subq, err := Query(v); err != nil {
					return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[should] failed to parse field").Cause(err)
				} else {
					boolQuery.AddShould(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := Query(vv.(map[string]interface{})); err != nil {
						return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[should] failed to parse field").Cause(err)
					} else {
						boolQuery.AddShould(subq)
					}
				}
			default:
				return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[bool] unsupported %s children type %T", k, v))
			}
		case "must":
			switch v := v.(type) {
			case map[string]interface{}:
				if subq, err := Query(v); err != nil {
					return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[must] failed to parse field").Cause(err)
				} else {
					boolQuery.AddMust(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := Query(vv.(map[string]interface{})); err != nil {
						return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[must] failed to parse field").Cause(err)
					} else {
						boolQuery.AddMust(subq)
					}
				}
			default:
				return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[bool] unsupported %s children type %T", k, v))
			}
		case "must_not":
			switch v := v.(type) {
			case map[string]interface{}:
				if subq, err := Query(v); err != nil {
					return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[must_not] failed to parse field").Cause(err)
				} else {
					boolQuery.AddMustNot(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := Query(vv.(map[string]interface{})); err != nil {
						return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[must_not] failed to parse field").Cause(err)
					} else {
						boolQuery.AddMustNot(subq)
					}
				}
			default:
				return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[bool] unsupported %s children type %T", k, v))
			}
		case "filter":
			filterQuery := bluge.NewBooleanQuery().SetBoost(0)
			switch v := v.(type) {
			case map[string]interface{}:
				if subq, err := Query(v); err != nil {
					return nil, meta.NewError(meta.ErrorTypeParsingException, "[filter] failed to parse field").Cause(err)
				} else {
					filterQuery.AddMust(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := Query(vv.(map[string]interface{})); err != nil {
						return nil, meta.NewError(meta.ErrorTypeParsingException, "[filter] failed to parse field").Cause(err)
					} else {
						filterQuery.AddMust(subq)
					}
				}
			default:
				return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[bool] unsupported %s children type %T", k, v))
			}
			boolQuery.AddMust(filterQuery)
		case "minimum_should_match":
			boolQuery.SetMinShould(int(v.(float64)))
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[bool] unsupported query type [%s]", k))
		}
	}

	return boolQuery, nil
}
