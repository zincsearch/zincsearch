package uquery

import (
	"github.com/blugelabs/bluge"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

func MultiPhraseQuery(iQuery v1.ZincQuery) (bluge.SearchRequest, error) {
	dateQuery := bluge.NewDateRangeQuery(iQuery.Query.StartTime, iQuery.Query.EndTime).SetField("@timestamp")

	var field string
	if iQuery.Query.Field != "" {
		field = iQuery.Query.Field
	} else {
		field = "_all"
	}

	multiPhraseQuery := bluge.NewMultiPhraseQuery(iQuery.Query.Terms).SetField(field)
	query := bluge.NewBooleanQuery().AddMust(dateQuery).AddMust(multiPhraseQuery)

	searchRequest := bluge.NewTopNSearch(iQuery.MaxResults, query).SortBy(iQuery.SortFields).WithStandardAggregations()

	return searchRequest, nil
}
