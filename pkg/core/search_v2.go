package core

import (
	"context"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/rs/zerolog/log"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	uquery "github.com/prabhatsharma/zinc/pkg/uquery/v2"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/parser/fields"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/parser/source"
)

func (index *Index) SearchV2(query *meta.ZincQuery) (*meta.SearchResponse, error) {
	searchRequest, err := uquery.ParseQueryDSL(query)
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

	var Hits []meta.Hit
	next, err := dmi.Next()
	for err == nil && next != nil {
		var id string
		var timestamp time.Time
		var sourceData map[string]interface{}
		var fieldsData map[string]interface{}
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			if field == "_source" {
				sourceData = source.Response(query.Source.(*meta.Source), value)
				if query.Fields != nil {
					fieldsData = fields.Response(query.Fields.([]*meta.Field), value)
				}
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

	return resp, nil
}
