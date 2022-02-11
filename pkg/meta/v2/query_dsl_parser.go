package v2

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"
	"github.com/prabhatsharma/zinc/pkg/startup"
)

// Parse parse query DSL and return searchRequest
func (q *ZincQuery) Parse() (bluge.SearchRequest, error) {
	if q.Size > int64(startup.MAX_RESULTS) {
		q.Size = int64(startup.MAX_RESULTS)
	}
	if q.Size == 0 {
		q.Size = 10
	}

	// parse track_total_hits

	// parse query
	query, err := q.parseQuery(q.Query)
	if err != nil {
		return nil, err
	}
	if query == nil {
		return nil, NewError(ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", q.Query))
	}

	request := bluge.NewTopNSearch(int(q.Size), query).WithStandardAggregations()

	// parse highlight

	// parse aggregations

	// parse fields

	// parse source

	// parse sort

	return request, nil
}

func (q *ZincQuery) parseQuery(query map[string]interface{}) (bluge.Query, error) {
	var cmd string
	var subq bluge.Query
	var err error
	for k, v := range query {
		if subq != nil && cmd != "" {
			return nil, NewError(ErrorTypeParsingException, fmt.Sprintf("[%s] malformed query, excepted [END_OBJECT] but found [FIELD_NAME] %s", cmd, k))
		}
		k := strings.ToLower(k)
		cmd = k
		v, ok := v.(map[string]interface{})
		if !ok {
			return nil, NewError(ErrorTypeParsingException, fmt.Sprintf("[%s] query doesn't support value type %T", k, v))
		}
		switch k {
		case "bool":
			if subq, err = q.parseBoolQuery(v); err != nil {
				return nil, NewError(ErrorTypeXContentParseException, "[bool] failed to parse field").Cause(err)
			}
		case "boosting":
			return nil, NewError(ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		case "match":
		case "match_bool_prefix":
		case "match_phrase":
		case "match_phrase_prefix":
		case "multi_match":
		case "match_all":
		case "match_none":
		case "combined_fields":
			return nil, NewError(ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		case "query_string":
		case "simple_query_string":
		case "exists":
		case "ids":
		case "range":
		case "prefix":
		case "fuzzy":
			if subq, err = q.parseFuzzyQuery(v); err != nil {
				return nil, NewError(ErrorTypeXContentParseException, "[fuzzy] failed to parse field").Cause(err)
			}
		case "wildcard":
			if subq, err = q.parseWildcardQuery(v); err != nil {
				return nil, NewError(ErrorTypeXContentParseException, "[wildcard] failed to parse field").Cause(err)
			}
		case "term":
		case "terms":
		case "terms_set":
			return nil, NewError(ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		case "geo_bounding_box":
			return nil, NewError(ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		case "geo_distance":
			return nil, NewError(ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		case "geo_polygon":
			return nil, NewError(ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		case "geo_shape":
			return nil, NewError(ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		default:
			return nil, NewError(ErrorTypeParsingException, fmt.Sprintf("[%s] query doesn't support", k))
		}
	}

	return subq, nil
}

func (q *ZincQuery) parseBoolQuery(query map[string]interface{}) (bluge.Query, error) {
	boolQuery := bluge.NewBooleanQuery()
	for k, v := range query {
		k := strings.ToLower(k)
		switch k {
		case "should":
			switch v := v.(type) {
			case map[string]interface{}:
				if subq, err := q.parseQuery(v); err != nil {
					return nil, NewError(ErrorTypeXContentParseException, "[should] failed to parse field").Cause(err)
				} else {
					boolQuery.AddShould(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := q.parseQuery(vv.(map[string]interface{})); err != nil {
						return nil, NewError(ErrorTypeXContentParseException, "[should] failed to parse field").Cause(err)
					} else {
						boolQuery.AddShould(subq)
					}
				}
			default:
				return nil, NewError(ErrorTypeParsingException, fmt.Sprintf("[bool] unsupported %s children type %T", k, v))
			}
		case "must":
			switch v := v.(type) {
			case map[string]interface{}:
				if subq, err := q.parseQuery(v); err != nil {
					return nil, NewError(ErrorTypeXContentParseException, "[must] failed to parse field").Cause(err)
				} else {
					boolQuery.AddMust(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := q.parseQuery(vv.(map[string]interface{})); err != nil {
						return nil, NewError(ErrorTypeXContentParseException, "[must] failed to parse field").Cause(err)
					} else {
						boolQuery.AddMust(subq)
					}
				}
			default:
				return nil, NewError(ErrorTypeParsingException, fmt.Sprintf("[bool] unsupported %s children type %T", k, v))
			}
		case "must_not":
			switch v := v.(type) {
			case map[string]interface{}:
				if subq, err := q.parseQuery(v); err != nil {
					return nil, NewError(ErrorTypeXContentParseException, "[must_not] failed to parse field").Cause(err)
				} else {
					boolQuery.AddMustNot(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := q.parseQuery(vv.(map[string]interface{})); err != nil {
						return nil, NewError(ErrorTypeXContentParseException, "[must_not] failed to parse field").Cause(err)
					} else {
						boolQuery.AddMustNot(subq)
					}
				}
			default:
				return nil, NewError(ErrorTypeParsingException, fmt.Sprintf("[bool] unsupported %s children type %T", k, v))
			}
		case "filter":
			filterQuery := bluge.NewBooleanQuery().SetBoost(0)
			switch v := v.(type) {
			case map[string]interface{}:
				if subq, err := q.parseQuery(v); err != nil {
					return nil, NewError(ErrorTypeParsingException, "[filter] failed to parse field").Cause(err)
				} else {
					filterQuery.AddMust(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := q.parseQuery(vv.(map[string]interface{})); err != nil {
						return nil, NewError(ErrorTypeParsingException, "[filter] failed to parse field").Cause(err)
					} else {
						filterQuery.AddMust(subq)
					}
				}
			default:
				return nil, NewError(ErrorTypeParsingException, fmt.Sprintf("[bool] unsupported %s children type %T", k, v))
			}
			boolQuery.AddMust(filterQuery)
		default:
			return nil, NewError(ErrorTypeParsingException, fmt.Sprintf("[bool] unsupported query type [%s]", k))
		}
	}

	return boolQuery, nil
}

func (q *ZincQuery) parseFuzzyQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, NewError(ErrorTypeParsingException, "[fuzzy] query doesn't support multiple fields")
	}

	field := ""
	value := new(FuzzyQuery)
	value.Boost = -1.0
	for k, v := range query {
		field = k
		switch v := v.(type) {
		case string:
			value.Value = v
		case map[string]interface{}:
			for k, v := range v {
				switch k {
				case "value":
					value.Value = v.(string)
				case "fuzziness":
					value.Fuzziness = v.(string)
				case "prefix_length":
					value.PrefixLength = v.(float64)
				case "boost":
					value.Boost = v.(float64)
				default:
					return nil, NewError(ErrorTypeParsingException, fmt.Sprintf("[fuzzy] unsupported children %s", k))
				}
			}
		default:
			return nil, NewError(ErrorTypeParsingException, fmt.Sprintf("[fuzzy] unsupported query type %s", k))
		}
	}

	subq := bluge.NewFuzzyQuery(value.Value).SetField(field)
	if value.Fuzziness != nil {
		switch v := value.Fuzziness.(type) {
		case string:
			// TODO: support other fuzziness: AUTO
		case float64:
			subq.SetFuzziness(int(v))
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

func (q *ZincQuery) parseWildcardQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, NewError(ErrorTypeParsingException, "[wildcard] query doesn't support multiple fields")
	}

	field := ""
	value := new(WildcardQuery)
	value.Boost = -1.0
	for k, v := range query {
		field = k
		switch v := v.(type) {
		case string:
			value.Value = v
		case map[string]interface{}:
			for k, v := range v {
				switch k {
				case "value":
					value.Value = v.(string)
				case "boost":
					value.Boost = v.(float64)
				default:
					return nil, NewError(ErrorTypeParsingException, fmt.Sprintf("[wildcard] unsupported children %s", k))
				}
			}
		default:
			return nil, NewError(ErrorTypeParsingException, fmt.Sprintf("[wildcard] unsupported query type %s", k))
		}
	}

	subq := bluge.NewWildcardQuery(value.Value).SetField(field)
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}
