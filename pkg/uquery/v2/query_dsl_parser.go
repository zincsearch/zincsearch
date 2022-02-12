package v2

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis/analyzer"
	querystr "github.com/blugelabs/query_string"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	"github.com/prabhatsharma/zinc/pkg/startup"
)

// ParseQueryDSL parse query DSL and return searchRequest
func ParseQueryDSL(q *meta.ZincQuery) (bluge.SearchRequest, error) {
	// parse size
	if q.Size > startup.MAX_RESULTS {
		q.Size = startup.MAX_RESULTS
	}
	if q.Size == 0 {
		q.Size = 10
	}

	// parse query
	query, err := parseQuery(q.Query)
	if err != nil {
		return nil, err
	}
	if query == nil {
		return nil, meta.NewError(meta.ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", q.Query))
	}

	// parse highlight

	// parse aggregations

	// create search request
	request := bluge.NewTopNSearch(q.Size, query).WithStandardAggregations()

	// parse from
	if q.From > 0 {
		request.SetFrom(q.From)
	}

	// parse explain
	if q.Explain {
		request.ExplainScores()
	}

	// parse fields

	// parse source

	// parse sort

	// parse track_total_hits
	// TODO support track_total_hits

	return request, nil
}

func parseQuery(query map[string]interface{}) (bluge.Query, error) {
	var subq bluge.Query
	var cmd string
	var err error
	for k, v := range query {
		if subq != nil && cmd != "" {
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[%s] malformed query, excepted [END_OBJECT] but found [FIELD_NAME] %s", cmd, k))
		}
		k := strings.ToLower(k)
		cmd = k
		v, ok := v.(map[string]interface{})
		if !ok {
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[%s] query doesn't support value type %T", k, v))
		}
		switch k {
		case "bool":
			if subq, err = parseBoolQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[bool] failed to parse field").Cause(err)
			}
		case "boosting":
			return nil, meta.NewError(meta.ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		case "match":
			if subq, err = parseMatchQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[match] failed to parse field").Cause(err)
			}
		case "match_bool_prefix":
			if subq, err = parseMatchBoolPrefixQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[match_bool_prefix] failed to parse field").Cause(err)
			}
		case "match_phrase":
			if subq, err = parseMatchPhraseQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[match_phrase] failed to parse field").Cause(err)
			}
		case "match_phrase_prefix":
			if subq, err = parseMatchPhrasePrefixQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[match_phrase_prefix] failed to parse field").Cause(err)
			}
		case "multi_match":
			if subq, err = parseMultiMatchQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[multi_match] failed to parse field").Cause(err)
			}
		case "match_all":
			if subq, err = parseMatchAllQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[match_all] failed to parse field").Cause(err)
			}
		case "match_none":
			if subq, err = parseMatchNoneQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[match_none] failed to parse field").Cause(err)
			}
		case "combined_fields":
			return nil, meta.NewError(meta.ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		case "query_string":
			if subq, err = parseQueryStringQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[query_string] failed to parse field").Cause(err)
			}
		case "simple_query_string":
			if subq, err = parseSimpleQueryStringQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[simple_query_string] failed to parse field").Cause(err)
			}
		case "exists":
			return nil, meta.NewError(meta.ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		case "ids":
			return nil, meta.NewError(meta.ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		case "range":
			if subq, err = parseRangeQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[range] failed to parse field").Cause(err)
			}
		case "prefix":
			if subq, err = parsePrefixQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[prefix] failed to parse field").Cause(err)
			}
		case "fuzzy":
			if subq, err = parseFuzzyQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[fuzzy] failed to parse field").Cause(err)
			}
		case "wildcard":
			if subq, err = parseWildcardQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[wildcard] failed to parse field").Cause(err)
			}
		case "term":
			if subq, err = parseTermQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[term] failed to parse field").Cause(err)
			}
		case "terms":
			if subq, err = parseTermsQuery(v); err != nil {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[terms] failed to parse field").Cause(err)
			}
		case "terms_set":
			return nil, meta.NewError(meta.ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		case "geo_bounding_box":
			return nil, meta.NewError(meta.ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		case "geo_distance":
			return nil, meta.NewError(meta.ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		case "geo_polygon":
			return nil, meta.NewError(meta.ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		case "geo_shape":
			return nil, meta.NewError(meta.ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", k))
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[%s] query doesn't support", k))
		}
	}

	return subq, nil
}

func parseBoolQuery(query map[string]interface{}) (bluge.Query, error) {
	boolQuery := bluge.NewBooleanQuery()
	for k, v := range query {
		k := strings.ToLower(k)
		switch k {
		case "should":
			switch v := v.(type) {
			case map[string]interface{}:
				if subq, err := parseQuery(v); err != nil {
					return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[should] failed to parse field").Cause(err)
				} else {
					boolQuery.AddShould(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := parseQuery(vv.(map[string]interface{})); err != nil {
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
				if subq, err := parseQuery(v); err != nil {
					return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[must] failed to parse field").Cause(err)
				} else {
					boolQuery.AddMust(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := parseQuery(vv.(map[string]interface{})); err != nil {
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
				if subq, err := parseQuery(v); err != nil {
					return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[must_not] failed to parse field").Cause(err)
				} else {
					boolQuery.AddMustNot(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := parseQuery(vv.(map[string]interface{})); err != nil {
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
				if subq, err := parseQuery(v); err != nil {
					return nil, meta.NewError(meta.ErrorTypeParsingException, "[filter] failed to parse field").Cause(err)
				} else {
					filterQuery.AddMust(subq)
				}
			case []interface{}:
				for _, vv := range v {
					if subq, err := parseQuery(vv.(map[string]interface{})); err != nil {
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

func parseMatchQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[match] query doesn't support multiple fields")
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
				switch k {
				case "query":
					value.Query = v.(string)
				case "analyzer":
					value.Analyzer = v.(string)
				case "operator":
					value.Operator = v.(string)
				case "fuzziness":
					value.Fuzziness = v.(string)
				case "prefix_length":
					value.PrefixLength = v.(float64)
				case "boost":
					value.Boost = v.(float64)
				default:
					return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[match] unsupported children %s", k))
				}
			}
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[match] unsupported query type %s", k))
		}
	}

	subq := bluge.NewMatchQuery(value.Query).SetField(field)
	if value.Analyzer != "" {
		switch value.Analyzer {
		case "standard":
			subq.SetAnalyzer(analyzer.NewStandardAnalyzer())
		default:
			// TODO: support analyzer
		}
	}
	if value.Operator != "" {
		op := strings.ToUpper(value.Operator)
		switch op {
		case "OR":
			subq.SetOperator(bluge.MatchQueryOperatorOr)
		case "AND":
			subq.SetOperator(bluge.MatchQueryOperatorAnd)
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[match] unsupported operator %s", op))
		}
	}
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

func parseMatchBoolPrefixQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[match_bool_prefix] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.MatchBoolPrefixQuery)
	for k, v := range query {
		field = k
		switch v := v.(type) {
		case string:
			value.Query = v
		case map[string]interface{}:
			for k, v := range v {
				switch k {
				case "query":
					value.Query = v.(string)
				case "analyzer":
					value.Analyzer = v.(string)
				default:
					return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[match_bool_prefix] unsupported children %s", k))
				}
			}
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[match_bool_prefix] unsupported query type %s", k))
		}
	}

	// TODO support analyzer
	zer := analyzer.NewStandardAnalyzer()
	if value.Analyzer != "" {
		switch value.Analyzer {
		case "standard":
			zer = analyzer.NewStandardAnalyzer()
		default:
			// TODO: support analyzer
		}
	}

	tokens := zer.Analyze([]byte(value.Query))
	subq := bluge.NewBooleanQuery()
	for i := 0; i < len(tokens); i++ {
		if i == len(tokens)-1 {
			subq.AddShould(bluge.NewPrefixQuery(string(tokens[i].Term)).SetField(field))
		} else {
			subq.AddShould(bluge.NewTermQuery(string(tokens[i].Term)).SetField(field))
		}
	}

	return subq, nil
}

func parseMatchPhraseQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[match_phrase] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.MatchPhraseQuery)
	value.Boost = -1.0
	for k, v := range query {
		field = k
		switch v := v.(type) {
		case string:
			value.Query = v
		case map[string]interface{}:
			for k, v := range v {
				switch k {
				case "query":
					value.Query = v.(string)
				case "analyzer":
					value.Analyzer = v.(string)
				case "boost":
					value.Boost = v.(float64)
				default:
					return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[match_phrase] unsupported children %s", k))
				}
			}
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[match_phrase] unsupported query type %s", k))
		}
	}

	subq := bluge.NewMatchPhraseQuery(value.Query).SetField(field)
	if value.Analyzer != "" {
		switch value.Analyzer {
		case "standard":
			subq.SetAnalyzer(analyzer.NewStandardAnalyzer())
		default:
			// TODO: support analyzer
		}
	}
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}

func parseMatchPhrasePrefixQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[match_phrase_prefix] query doesn't support multiple fields")
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
				switch k {
				case "query":
					value.Query = v.(string)
				case "analyzer":
					value.Analyzer = v.(string)
				case "boost":
					value.Boost = v.(float64)
				default:
					return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[match_phrase_prefix] unsupported children %s", k))
				}
			}
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[match_phrase_prefix] unsupported query type %s", k))
		}
	}

	// TODO support analyzer
	zer := analyzer.NewStandardAnalyzer()
	if value.Analyzer != "" {
		switch value.Analyzer {
		case "standard":
			zer = analyzer.NewStandardAnalyzer()
		default:
			// TODO: support analyzer
		}
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

func parseMultiMatchQuery(query map[string]interface{}) (bluge.Query, error) {
	value := new(meta.MultiMatchQuery)
	for k, v := range query {
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
		case "type":
			value.Type = v.(string)
		case "operator":
			value.Operator = v.(string)
		case "minimum_should_match":
			value.MinimumShouldMatch = v.(float64)
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[multi_match] unsupported children %s", k))
		}
	}

	// TODO support analyzer
	zer := analyzer.NewStandardAnalyzer()
	if value.Analyzer != "" {
		switch value.Analyzer {
		case "standard":
			zer = analyzer.NewStandardAnalyzer()
		default:
			// TODO: support analyzer
		}
	}

	var operator bluge.MatchQueryOperator = bluge.MatchQueryOperatorOr
	if value.Operator != "" {
		op := strings.ToUpper(value.Operator)
		switch op {
		case "OR":
			operator = bluge.MatchQueryOperatorOr
		case "AND":
			operator = bluge.MatchQueryOperatorAnd
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[multi_match] unsupported operator %s", op))
		}
	}

	subq := bluge.NewBooleanQuery()
	if value.MinimumShouldMatch > 0 {
		subq.SetMinShould(int(value.MinimumShouldMatch))
	}
	for _, field := range value.Fields {
		subq.AddShould(bluge.NewMatchQuery(value.Query).SetField(field).SetAnalyzer(zer).SetOperator(operator))
	}

	return subq, nil
}

func parseMatchAllQuery(query map[string]interface{}) (bluge.Query, error) {
	return bluge.NewMatchAllQuery(), nil
}

func parseMatchNoneQuery(query map[string]interface{}) (bluge.Query, error) {
	return bluge.NewMatchNoneQuery(), nil
}

func parseQueryStringQuery(query map[string]interface{}) (bluge.Query, error) {
	value := new(meta.QueryStringQuery)
	for k, v := range query {
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
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[query_string] unsupported children %s", k))
		}
	}

	// TODO fields
	// TODO default_field
	// TODO default_operator
	// TODO boost

	// TODO support analyzer
	zer := analyzer.NewStandardAnalyzer()
	if value.Analyzer != "" {
		switch value.Analyzer {
		case "standard":
			zer = analyzer.NewStandardAnalyzer()
		default:
			// TODO: support analyzer
		}
	}

	options := querystr.DefaultOptions()
	options.WithDefaultAnalyzer(zer)
	return querystr.ParseQueryString(value.Query, options)
}

func parseSimpleQueryStringQuery(query map[string]interface{}) (bluge.Query, error) {
	return parseQueryStringQuery(query)
}

func parseRangeQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[range] query doesn't support multiple fields")
	}

	return nil, nil
}

func parsePrefixQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[prefix] query doesn't support multiple fields")
	}

	return nil, nil
}

func parseFuzzyQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[fuzzy] query doesn't support multiple fields")
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
					return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[fuzzy] unsupported children %s", k))
				}
			}
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[fuzzy] unsupported query type %s", k))
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

func parseWildcardQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[wildcard] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.WildcardQuery)
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
					return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[wildcard] unsupported children %s", k))
				}
			}
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[wildcard] unsupported query type %s", k))
		}
	}

	subq := bluge.NewWildcardQuery(value.Value).SetField(field)
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}

func parseTermQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[term] query doesn't support multiple fields")
	}

	return nil, nil
}

func parseTermsQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[terms] query doesn't support multiple fields")
	}

	return nil, nil
}
