package v2

import (
	"fmt"

	"github.com/blugelabs/bluge"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	"github.com/prabhatsharma/zinc/pkg/startup"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/parser/query"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/parser/source"
)

// ParseQueryDSL parse query DSL and return searchRequest
func ParseQueryDSL(q *meta.ZincQuery) (bluge.SearchRequest, error) {
	// parse size
	if q.Size == 0 {
		q.Size = 10
	}
	if q.Size > startup.MAX_RESULTS {
		q.Size = startup.MAX_RESULTS
	}

	// parse query
	query, err := query.Query(q.Query)
	if err != nil {
		return nil, err
	}
	if query == nil {
		return nil, meta.NewError(meta.ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", q.Query))
	}

	// parse highlight

	// parse aggregations

	// create search request
	request := bluge.NewTopNSearch(q.Size, query).WithStandardAggregations()

	// parse from
	if q.From > 0 {
		request.SetFrom(q.From)
	}

	// parse explain
	if q.Explain {
		request.ExplainScores()
	}

	// parse fields

	// parse source
	if q.Source != nil {
		if q.Source, err = source.Request(q.Source); err != nil {
			return nil, err
		}
	}

	// parse sort

	// parse track_total_hits
	// TODO support track_total_hits

	return request, nil
}
