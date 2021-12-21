package uquery

import (
	"github.com/blugelabs/bluge"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

func DateRangeQuery(iQuery v1.ZincQuery) (bluge.SearchRequest, error) {
	dateQuery := bluge.NewDateRangeQuery(iQuery.Query.StartTime, iQuery.Query.EndTime).SetField("@timestamp")
	query := bluge.NewBooleanQuery().AddMust(dateQuery)

	searchRequest := bluge.NewTopNSearch(iQuery.MaxResults, query).SortBy(iQuery.SortFields).WithStandardAggregations()

	return searchRequest, nil
}
