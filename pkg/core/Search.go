package core

import (
	"context"
	"encoding/json"
	"time"

	"github.com/blugelabs/bluge"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
	"github.com/prabhatsharma/zinc/pkg/uquery"
	"github.com/rs/zerolog/log"
)

func (index *Index) Search(iQuery v1.ZincQuery) (v1.SearchResponse, error) {
	var Hits []v1.Hit

	var searchRequest bluge.SearchRequest

	if iQuery.MaxResults == 0 {
		iQuery.MaxResults = 20
	}

	var err error

	switch iQuery.SearchType {
	case "alldocuments":
		searchRequest, err = uquery.AllDocuments(iQuery)
	case "wildcard":
		searchRequest, err = uquery.WildcardQuery(iQuery)
	case "fuzzy":
		searchRequest, err = uquery.FuzzyQuery(iQuery)
	case "term":
		searchRequest, err = uquery.TermQuery(iQuery)
	case "daterange":
		searchRequest, err = uquery.DateRangeQuery(iQuery)
	case "matchall":
		searchRequest, err = uquery.MatchAllQuery(iQuery)
	case "match":
		searchRequest, err = uquery.MatchQuery(iQuery)
	case "matchphrase":
		searchRequest, err = uquery.MatchPhraseQuery(iQuery)
	case "multiphrase":
		searchRequest, err = uquery.MultiPhraseQuery(iQuery)
	case "prefix":
		searchRequest, err = uquery.PrefixQuery(iQuery)
	case "querystring":
		searchRequest, err = uquery.QueryStringQuery(iQuery)
	}

	if err != nil {
		resp := v1.SearchResponse{
			Error: err.Error(),
		}

		return resp, err
	}

	// sample time range aggregation start
	// timestampAggregation := aggregations.DateRanges(search.Field("@timestamp"))
	// daterange1 := aggregations.NewDateRange(time.Now().Add(-time.Hour*24*30), time.Now())
	// timestampAggregation.AddRange(daterange1)
	// searchRequest.AddAggregation("@timestamp", timestampAggregation)
	// sample time range aggregation end

	writer := index.Writer

	reader, err := writer.Reader()
	if err != nil {
		log.Printf("error accessing reader: %v", err)
	}

	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Printf("error executing search: %v", err)
	}

	// highlighter := highlight.NewANSIHighlighter()

	// iterationStartTime := time.Now()
	next, err := dmi.Next()
	for err == nil && next != nil {
		var result map[string]interface{}
		var id string
		var timestamp time.Time
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			if field == "_source" {
				json.Unmarshal(value, &result)
				return true
			} else if field == "_id" {
				id = string(value)
				return true
			} else if field == "@timestamp" {
				timestamp, _ = bluge.DecodeDateTime(value)
				return true
			}
			return true
		})
		if err != nil {
			log.Printf("error accessing stored fields: %v", err)
		}

		hit := v1.Hit{
			Index:     index.Name,
			Type:      index.Name,
			ID:        id,
			Score:     next.Score,
			Timestamp: timestamp,
			Source:    result,
		}

		next, err = dmi.Next()
		// results = append(results, result)

		Hits = append(Hits, hit)
	}
	if err != nil {
		log.Printf("error iterating results: %v", err)
	}

	// fmt.Println("Got results after data load from disk in: ", time.Since(iterationStartTime))
	resp := v1.SearchResponse{
		// Took: int(time.Since(searchStart).Milliseconds()),
		Took:     int(dmi.Aggregations().Duration().Milliseconds()),
		MaxScore: dmi.Aggregations().Metric("max_score"),
		// Buckets:  dmi.Aggregations().Buckets("@timestamp"),
		Hits: v1.Hits{
			Total: v1.Total{
				Value: int(dmi.Aggregations().Count()),
			},
			Hits: Hits,
		},
	}

	reader.Close()

	return resp, nil
}
