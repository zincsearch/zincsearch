package uquery

import (
	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis/analyzer"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

func MatchQuery(iQuery v1.ZincQuery) (bluge.SearchRequest, error) {
	// startTime := time.Now()

	// dateQuery := bluge.NewDateRangeQuery(iQuery.Query.StartTime, iQuery.Query.EndTime)
	dateQuery := bluge.NewDateRangeQuery(iQuery.Query.StartTime, iQuery.Query.EndTime).SetField("@timestamp")
	dateQuery.SetField("@timestamp")

	// log.Println("Start time", iQuery.Query.StartTime, ", End time: ", iQuery.Query.EndTime)

	var field string
	if iQuery.Query.Field != "" {
		field = iQuery.Query.Field
	} else {
		field = "_all"
	}

	matchQuery := bluge.NewMatchQuery(iQuery.Query.Term).SetField(field).SetAnalyzer(analyzer.NewStandardAnalyzer())
	query := bluge.NewBooleanQuery().AddMust(dateQuery).AddMust(matchQuery)

	searchRequest := buildRequest(iQuery, query)

	// endTime := time.Now()
	// log.Println("Query time: ", endTime.Sub(startTime))

	return searchRequest, nil
}
