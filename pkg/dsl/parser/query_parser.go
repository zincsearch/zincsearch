package parser

import (
	"fmt"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/search"

	"github.com/prabhatsharma/zinc/pkg/dsl/meta"
	"github.com/prabhatsharma/zinc/pkg/dsl/parser/aggregation"
	"github.com/prabhatsharma/zinc/pkg/dsl/parser/fields"
	"github.com/prabhatsharma/zinc/pkg/dsl/parser/query"
	"github.com/prabhatsharma/zinc/pkg/dsl/parser/sort"
	"github.com/prabhatsharma/zinc/pkg/dsl/parser/source"
	"github.com/prabhatsharma/zinc/pkg/startup"
)

// ParseQuery parse query DSL and return searchRequest
func ParseQuery(q *meta.ZincQuery, mappings *meta.Mappings) (bluge.SearchRequest, error) {
	// parse size
	if q.Size == 0 {
		q.Size = 10
	}
	if q.Size > startup.LoadMaxResults() {
		q.Size = startup.LoadMaxResults()
	}

	// parse query
	query, err := query.Query(q.Query)
	if err != nil {
		return nil, err
	}
	if query == nil {
		return nil, meta.NewError(meta.ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", q.Query))
	}

	// create search request
	request := bluge.NewTopNSearch(q.Size, query).WithStandardAggregations()

	// parse highlight
	// TODO: highlight
	if q.Highlight != nil {
		request.IncludeLocations()
	}

	// parse from
	if q.From > 0 {
		request.SetFrom(q.From)
	}

	// parse explain
	if q.Explain {
		request.ExplainScores()
	}

	// parse aggregations
	if q.Aggregations != nil {
		if err := aggregation.Request(request, q.Aggregations, mappings); err != nil {
			return nil, err
		}
	}

	// parse fields
	if q.Fields != nil {
		if v, ok := q.Fields.([]interface{}); ok {
			if q.Fields, err = fields.Request(v); err != nil {
				return nil, err
			}
		}
	}

	// parse source
	if q.Source, err = source.Request(q.Source); err != nil {
		return nil, err
	}

	// parse sort
	if q.Sort != nil {
		if q.Sort, err = sort.Request(q.Sort); err != nil {
			return nil, err
		}
		if q.Sort != nil {
			request.SortByCustom(q.Sort.(search.SortOrder))
		}
	}

	// pagenation
	// TODO: search after PIT support

	return request, nil
}
