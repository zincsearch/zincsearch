package uquery

import (
	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis/analyzer"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

func MatchPhraseQuery(iQuery *v1.ZincQuery) (bluge.SearchRequest, error) {
	dateQuery := bluge.NewDateRangeQuery(iQuery.Query.StartTime, iQuery.Query.EndTime).SetField("@timestamp")

	var field string
	if iQuery.Query.Field != "" {
		field = iQuery.Query.Field
	} else {
		field = "_all"
	}

	matchPhraseQuery := bluge.NewMatchPhraseQuery(iQuery.Query.Term).SetField(field).SetAnalyzer(analyzer.NewStandardAnalyzer())
	query := bluge.NewBooleanQuery().AddMust(dateQuery).AddMust(matchPhraseQuery)

	searchRequest := buildRequest(iQuery, query)

	return searchRequest, nil
}
