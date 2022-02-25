package core

import (
	"context"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/search"
	"github.com/blugelabs/bluge/search/highlight"
	"github.com/rs/zerolog/log"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	parser "github.com/prabhatsharma/zinc/pkg/uquery/v2"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/fields"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/source"
)

func (index *Index) SearchV2(query *meta.ZincQuery) (*meta.SearchResponse, error) {
	searchRequest, err := parser.ParseQueryDSL(query, index.CachedMappings, index.CachedAnalyzers)
	if err != nil {
		return nil, err
	}

	reader, err := index.Writer.Reader()
	if err != nil {
		log.Printf("index.SearchV2: error accessing reader: %v", err)
		return nil, err
	}
	defer reader.Close()

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
			return &meta.SearchResponse{
				TimedOut: true,
				Error:    err.Error(),
				Hits:     meta.Hits{Hits: []meta.Hit{}},
			}, nil
		}
		return nil, err
	}

	return searchV2(dmi, query, index.CachedMappings)
}

func searchV2(dmi search.DocumentMatchIterator, query *meta.ZincQuery, mappings *meta.Mappings) (*meta.SearchResponse, error) {
	resp := &meta.SearchResponse{
		Hits: meta.Hits{Hits: []meta.Hit{}},
	}

	// highlight
	var highlighter *highlight.SimpleHighlighter
	if query.Highlight != nil {
		if len(query.Highlight.PreTags) > 0 && len(query.Highlight.PostTags) > 0 {
			highlighter = highlight.NewHTMLHighlighterTags(query.Highlight.PreTags[0], query.Highlight.PostTags[0])
		} else {
			highlighter = highlight.NewHTMLHighlighter()
		}
	}

	Hits := make([]meta.Hit, 0)
	next, err := dmi.Next()
	for err == nil && next != nil {
		var id string
		var indexName string
		var timestamp time.Time
		var sourceData map[string]interface{}
		var fieldsData map[string]interface{}
		var highlightData map[string]interface{}
		if query.Highlight != nil {
			highlightData = make(map[string]interface{})
		}
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			switch field {
			case "_id":
				id = string(value)
			case "_index":
				indexName = string(value)
			case "@timestamp":
				timestamp, _ = bluge.DecodeDateTime(value)
			case "_source":
				sourceData = source.Response(query.Source.(*meta.Source), value)
				if query.Fields != nil {
					fieldsData = fields.Response(query.Fields.([]*meta.Field), value, mappings)
				}
			default:
				// highlight
				if query.Highlight != nil && query.Highlight.Fields != nil {
					if options, ok := query.Highlight.Fields[field]; ok {
						if v, ok := next.Locations[field]; ok {
							if len(options.PreTags) > 0 && len(options.PostTags) > 0 {
								highlighter := highlight.NewHTMLHighlighterTags(options.PreTags[0], options.PostTags[0])
								highlightData[field] = highlighter.BestFragments(v, value, options.NumberOfFragments)
							} else {
								highlightData[field] = highlighter.BestFragments(v, value, options.NumberOfFragments)
							}
						}
					}
				}
			}

			return true
		})
		if err != nil {
			log.Printf("core.SearchV2: error accessing stored fields: %v", err)
			continue
		}

		hit := meta.Hit{
			Index:     indexName,
			Type:      "_doc",
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
		log.Printf("core.SearchV2: error iterating results: %v", err)
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

	if err := parser.FormatResponse(resp, query, dmi.Aggregations()); err != nil {
		log.Printf("core.SearchV2: error format response: %v", err)
	}

	return resp, nil
}
