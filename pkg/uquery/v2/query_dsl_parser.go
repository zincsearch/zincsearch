package v2

import (
	"fmt"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/search"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	"github.com/prabhatsharma/zinc/pkg/startup"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/parser/aggregation"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/parser/fields"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/parser/query"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/parser/sort"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/parser/source"
)

// ParseQueryDSL parse query DSL and return searchRequest
func ParseQueryDSL(q *meta.ZincQuery, mappings map[string]string) (bluge.SearchRequest, error) {
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

	// parse highlight
	// TODO: highlight

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
