package parser

import (
	"fmt"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/search"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	"github.com/prabhatsharma/zinc/pkg/startup"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/aggregation"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/fields"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/highlight"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/query"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/sort"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/source"
)

// ParseQueryDSL parse query DSL and return searchRequest
func ParseQueryDSL(q *meta.ZincQuery, mappings *meta.Mappings, analyzers map[string]*analysis.Analyzer) (bluge.SearchRequest, error) {
	// parse size
	if q.Size == 0 {
		q.Size = 10
	}
	if q.Size > startup.LoadMaxResults() {
		q.Size = startup.LoadMaxResults()
	}

	// parse query
	query, err := query.Query(q.Query, mappings, analyzers)
	if err != nil {
		return nil, err
	}
	if query == nil {
		return nil, errors.New(errors.ErrorTypeNotImplemented, fmt.Sprintf("[%s] query doesn't support", q.Query))
	}

	// create search request
	request := bluge.NewTopNSearch(q.Size, query).WithStandardAggregations()

	// parse highlight
	if q.Highlight != nil {
		highlight.Request(q.Highlight)
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
