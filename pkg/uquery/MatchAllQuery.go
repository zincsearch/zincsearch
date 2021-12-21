package uquery

import (
	"github.com/blugelabs/bluge"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

func MatchAllQuery(iQuery v1.ZincQuery) (bluge.SearchRequest, error) {
	dateQuery := bluge.NewDateRangeQuery(iQuery.Query.StartTime, iQuery.Query.EndTime).SetField("@timestamp")

	var field string
	if iQuery.Query.Field != "" {
		field = iQuery.Query.Field
	} else {
		field = "_all"
	}

	fuzzyQuery := bluge.NewFuzzyQuery(iQuery.Query.Term).SetField(field)
	query := bluge.NewBooleanQuery().AddMust(dateQuery).AddMust(fuzzyQuery)

	searchRequest := bluge.NewTopNSearch(iQuery.MaxResults, query).SortBy(iQuery.SortFields).WithStandardAggregations()

	return searchRequest, nil
}
