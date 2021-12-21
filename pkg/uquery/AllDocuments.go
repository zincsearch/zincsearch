package uquery

import (
	"github.com/blugelabs/bluge"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

func AllDocuments(iQuery v1.ZincQuery) (bluge.SearchRequest, error) {

	dateQuery := bluge.NewDateRangeQuery(iQuery.Query.StartTime, iQuery.Query.EndTime).SetField("@timestamp")

	allquery := bluge.NewMatchAllQuery()

	query := bluge.NewBooleanQuery().AddMust(dateQuery).AddMust(allquery)

	searchRequest := bluge.NewTopNSearch(20, query).SortBy(iQuery.SortFields)

	return searchRequest, nil

}
