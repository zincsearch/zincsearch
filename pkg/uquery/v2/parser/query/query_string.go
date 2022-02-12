package query

import (
	"fmt"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis/analyzer"
	querystr "github.com/blugelabs/query_string"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func QueryStringQuery(query map[string]interface{}) (bluge.Query, error) {
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
