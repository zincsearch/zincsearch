package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	zincanalysis "github.com/prabhatsharma/zinc/pkg/uquery/v2/analysis"
)

func MultiMatchQuery(query map[string]interface{}, mappings *meta.Mappings, analyzers map[string]*analysis.Analyzer) (bluge.Query, error) {
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
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[multi_match] unknown field [%s]", k))
		}
	}

	var zer *analysis.Analyzer
	if value.Analyzer != "" {
		zer, _ = zincanalysis.QueryAnalyzer(analyzers, value.Analyzer)
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
			return nil, errors.New(errors.ErrorTypeIllegalArgumentException, fmt.Sprintf("[multi_match] unknown operator %s", op))
		}
	}

	subq := bluge.NewBooleanQuery()
	if value.MinimumShouldMatch > 0 {
		subq.SetMinShould(int(value.MinimumShouldMatch))
	}
	for _, field := range value.Fields {
		subqq := bluge.NewMatchQuery(value.Query).SetField(field).SetOperator(operator)
		if zer != nil {
			subqq.SetAnalyzer(zer)
		}
		subq.AddShould(subqq)
	}

	return subq, nil
}
