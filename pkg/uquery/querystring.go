package uquery

import (
	"fmt"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis/analyzer"
	querystr "github.com/blugelabs/query_string"
	v1 "github.com/zinclabs/zinc/pkg/meta/v1"
)

func QueryStringQuery(iQuery *v1.ZincQuery) (bluge.SearchRequest, error) {
	options := querystr.DefaultOptions()
	options.WithDefaultAnalyzer(analyzer.NewStandardAnalyzer())
	userQuery, err := querystr.ParseQueryString(iQuery.Query.Term, options)
	if err != nil {
		return nil, fmt.Errorf("error parsing query string '%s': %s", iQuery.Query.Term, err.Error())
	}

	dateQuery := bluge.NewDateRangeQuery(iQuery.Query.StartTime, iQuery.Query.EndTime).SetField("@timestamp")
	finalQuery := bluge.NewBooleanQuery().AddMust(dateQuery).AddMust(userQuery)

	// sortFields := []string{"-@timestamp"} // adding a - (minus) before the field name will sort the field in descending order

	searchRequest := buildRequest(iQuery, finalQuery)

	return searchRequest, nil

}
