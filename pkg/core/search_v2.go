package core

import (
	"context"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/search/highlight"
	"github.com/rs/zerolog/log"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	uquery "github.com/prabhatsharma/zinc/pkg/uquery/v2"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/parser/fields"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/parser/source"
)

func (index *Index) SearchV2(query *meta.ZincQuery) (*meta.SearchResponse, error) {
	mappings, _ := index.GetStoredMappings()
	searchRequest, err := uquery.ParseQueryDSL(query, mappings)
	if err != nil {
		return nil, err
	}

	reader, err := index.Writer.Reader()
	if err != nil {
		log.Printf("index.SearchV2: error accessing reader: %v", err)
		return nil, err
	}
	defer reader.Close()

	resp := new(meta.SearchResponse)
	ctx := context.Background()
	var cancel context.CancelFunc
	if query.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(query.Timeout)*time.Second)
		defer cancel()
	}
	dmi, err := reader.Search(ctx, searchRequest)
	if err != nil {
		log.Printf("index.SearchV2: error executing search: %v", err)
		if err == context.DeadlineExceeded {
			resp.TimedOut = true
			resp.Error = err.Error()
			return resp, err
		}
		return nil, err
	}

	higher := highlight.NewHTMLHighlighter()

	var Hits []meta.Hit
	next, err := dmi.Next()
	for err == nil && next != nil {
		var id string
		var timestamp time.Time
		var sourceData map[string]interface{}
		var fieldsData map[string]interface{}
		var highlightData map[string]interface{}
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			switch field {
			case "_id":
				id = string(value)
			case "@timestamp":
				timestamp, _ = bluge.DecodeDateTime(value)
			case "_source":
				sourceData = source.Response(query.Source.(*meta.Source), value)
				if query.Fields != nil {
					fieldsData = fields.Response(query.Fields.([]*meta.Field), value)
				}
			default:
				// highlight
				if query.Highlight != nil && query.Highlight.Fields != nil {
					if highlightData == nil {
						highlightData = make(map[string]interface{})
					}
					// TODO support highlight options
					if options, ok := query.Highlight.Fields[field]; ok {
						if v, ok := next.Locations[field]; ok {
							options.NumberOfFragments = 1 // TODO support multiple fragments
							highlightData[field] = higher.BestFragments(v, value, options.NumberOfFragments)
						}
					}
				}
			}

			return true
		})
		if err != nil {
			log.Printf("index.SearchV2: error accessing stored fields: %v", err)
		}

		hit := meta.Hit{
			Index:     index.Name,
			Type:      index.Name,
			ID:        id,
			Score:     next.Score,
			Timestamp: timestamp,
			Source:    sourceData,
			Fields:    fieldsData,
			Highlight: highlightData,
		}
		Hits = append(Hits, hit)

		next, err = dmi.Next()
	}
	if err != nil {
		log.Printf("index.SearchV2: error iterating results: %v", err)
	}

	resp.Took = int(dmi.Aggregations().Duration().Milliseconds())
	resp.Shards = meta.Shards{Total: 1, Successful: 1}
	resp.Hits = meta.Hits{
		Total: meta.Total{
			Value: int(dmi.Aggregations().Count()),
		},
		MaxScore: dmi.Aggregations().Metric("max_score"),
		Hits:     Hits,
	}

	if err := uquery.FormatResponse(resp, query, dmi.Aggregations()); err != nil {
		log.Printf("index.SearchV2: error format response: %v", err)
	}

	return resp, nil
}
