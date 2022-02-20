package core

import (
	"context"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/rs/zerolog/log"

	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
	"github.com/prabhatsharma/zinc/pkg/startup"
	"github.com/prabhatsharma/zinc/pkg/uquery"
)

func (index *Index) Search(iQuery *v1.ZincQuery) (*v1.SearchResponse, error) {
	var searchRequest bluge.SearchRequest
	if iQuery.MaxResults > startup.LoadMaxResults() {
		iQuery.MaxResults = startup.LoadMaxResults()
	}

	sourceCtl := &v1.Source{Enable: true}
	switch iQuery.Source.(type) {
	case bool:
		sourceCtl.Enable = iQuery.Source.(bool)
	case []interface{}:
		v := iQuery.Source.([]interface{})
		sourceCtl.Fields = make(map[string]bool, len(v))
		for _, field := range v {
			if fv, ok := field.(string); ok {
				sourceCtl.Fields[fv] = true
			}
		}
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
	default:
		// default use alldocuments search
		searchRequest, err = uquery.AllDocuments(iQuery)
	}

	if err != nil {
		return &v1.SearchResponse{
			Error: err.Error(),
		}, err
	}

	// handle aggregations
	err = uquery.AddAggregations(searchRequest, iQuery.Aggregations, index.CachedMappings)
	if err != nil {
		return &v1.SearchResponse{
			Error: err.Error(),
		}, err
	}

	reader, err := index.Writer.Reader()
	if err != nil {
		log.Printf("error accessing reader: %v", err)
	}
	defer reader.Close()

	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Printf("error executing search: %v", err)
	}

	var Hits []v1.Hit

	next, err := dmi.Next()
	for err == nil && next != nil {
		var result map[string]interface{}
		var id string
		var timestamp time.Time
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			switch field {
			case "_id":
				id = string(value)
			case "@timestamp":
				timestamp, _ = bluge.DecodeDateTime(value)
			case "_source":
				result = uquery.HandleSource(sourceCtl, value)
			default:
			}

			return true
		})
		if err != nil {
			log.Printf("error accessing stored fields: %v", err)
		}

		hit := v1.Hit{
			Index:     index.Name,
			Type:      "_doc",
			ID:        id,
			Score:     next.Score,
			Timestamp: timestamp,
			Source:    result,
		}
		Hits = append(Hits, hit)

		next, err = dmi.Next()
	}
	if err != nil {
		log.Printf("error iterating results: %v", err)
	}

	resp := &v1.SearchResponse{
		Took: int(dmi.Aggregations().Duration().Milliseconds()),
		Hits: v1.Hits{
			Total: v1.Total{
				Value: int(dmi.Aggregations().Count()),
			},
			MaxScore: dmi.Aggregations().Metric("max_score"),
			Hits:     Hits,
		},
	}

	if len(iQuery.Aggregations) > 0 {
		resp.Aggregations, err = uquery.ParseAggregations(dmi.Aggregations())
		if err != nil {
			log.Printf("error parse aggregation results: %v", err)
		}
		if len(resp.Aggregations) > 0 {
			delete(resp.Aggregations, "count")
			delete(resp.Aggregations, "duration")
			delete(resp.Aggregations, "max_score")
		}
	}

	return resp, nil
}
