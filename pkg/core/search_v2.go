package core

import (
	"context"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/rs/zerolog/log"

	v2 "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func (index *Index) SearchV2(query *v2.ZincQuery) (*v2.SearchResponse, error) {
	resp := new(v2.SearchResponse)
	sourceCtl := &v2.Source{Enable: true}
	switch query.Source.(type) {
	case bool:
		sourceCtl.Enable = query.Source.(bool)
	case []interface{}:
		v := query.Source.([]interface{})
		sourceCtl.Fields = make(map[string]bool, len(v))
		for _, field := range v {
			if fv, ok := field.(string); ok {
				sourceCtl.Fields[fv] = true
			}
		}
	}

	searchRequest, err := query.Parse()
	if err != nil {
		return nil, err
	}

	reader, err := index.Writer.Reader()
	if err != nil {
		log.Printf("error accessing reader: %v", err)
		return nil, err
	}
	defer reader.Close()

	ctx := context.Background()
	var cancel context.CancelFunc
	if query.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(query.Timeout)*time.Millisecond)
		defer cancel()
	}
	dmi, err := reader.Search(ctx, searchRequest)
	if err != nil {
		log.Printf("error executing search: %v", err)
		if err == context.DeadlineExceeded {
			resp.TimedOut = true
			resp.Error = err.Error()
			return resp, err
		}
		return nil, err
	}

	var Hits []v2.Hit
	next, err := dmi.Next()
	for err == nil && next != nil {
		var result map[string]interface{}
		var id string
		var timestamp time.Time
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			if field == "_source" {
				result = v2.HandleSource(sourceCtl, value)
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

		hit := v2.Hit{
			Index:     index.Name,
			Type:      index.Name,
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

	resp.Took = int(dmi.Aggregations().Duration().Milliseconds())
	resp.Hits = v2.Hits{
		Total: v2.Total{
			Value: int(dmi.Aggregations().Count()),
		},
		MaxScore: dmi.Aggregations().Metric("max_score"),
		Hits:     Hits,
	}

	return resp, nil
}
