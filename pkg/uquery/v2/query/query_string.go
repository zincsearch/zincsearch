package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/analyzer"
	querystr "github.com/blugelabs/query_string"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	zincanalysis "github.com/prabhatsharma/zinc/pkg/uquery/v2/analysis"
)

func QueryStringQuery(query map[string]interface{}, mappings *meta.Mappings, analyzers map[string]*analysis.Analyzer) (bluge.Query, error) {
	value := new(meta.QueryStringQuery)
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
		case "default_field":
			value.DefaultField = v.(string)
		case "default_operator":
			value.DefaultOperator = v.(string)
		case "boost":
			value.Boost = v.(float64)
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[query_string] unsupported children %s", k))
		}
	}

	// TODO fields
	// TODO default_field
	// TODO default_operator
	// TODO boost

	var zer *analysis.Analyzer
	if value.Analyzer != "" {
		zer, _ = zincanalysis.QueryAnalyzer(analyzers, value.Analyzer)
	}
	if zer == nil {
		zer = analyzer.NewStandardAnalyzer()
	}

	options := querystr.DefaultOptions()
	options.WithDefaultAnalyzer(zer)
	return querystr.ParseQueryString(value.Query, options)
}
