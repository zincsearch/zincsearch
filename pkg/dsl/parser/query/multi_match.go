package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis/analyzer"

	"github.com/prabhatsharma/zinc/pkg/dsl/meta"
)

func MultiMatchQuery(query map[string]interface{}) (bluge.Query, error) {
	value := new(meta.MultiMatchQuery)
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
		case "type":
			value.Type = v.(string)
		case "operator":
			value.Operator = v.(string)
		case "minimum_should_match":
			value.MinimumShouldMatch = v.(float64)
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[multi_match] unknown field [%s]", k))
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
			return nil, meta.NewError(meta.ErrorTypeIllegalArgumentException, fmt.Sprintf("[multi_match] unknown operator %s", op))
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
