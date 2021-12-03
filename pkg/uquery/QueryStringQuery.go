package uquery

import (
	"fmt"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis/analyzer"
	querystr "github.com/blugelabs/query_string"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

func QueryStringQuery(iQuery v1.ZincQuery) (bluge.SearchRequest, error) {
	options := querystr.DefaultOptions()
	options.WithDefaultAnalyzer(analyzer.NewStandardAnalyzer())
	userQuery, err := querystr.ParseQueryString(iQuery.Query.Term, options)
	if err != nil {
		return nil, fmt.Errorf("error parsing query string '%s': %v", iQuery.Query.Term, err)
	}

	dateQuery := bluge.NewDateRangeQuery(iQuery.Query.StartTime, iQuery.Query.EndTime).SetField("@timestamp")
	finalQuery := bluge.NewBooleanQuery().AddMust(dateQuery).AddMust(userQuery)

	searchRequest := bluge.NewTopNSearch(1000, finalQuery).WithStandardAggregations()

	return searchRequest, nil

}
