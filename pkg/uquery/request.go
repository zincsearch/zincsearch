package uquery

import (
	"github.com/blugelabs/bluge"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

// buildRequest combines the ZincQuery with the bluge Query to create a SearchRequest
func buildRequest(iQuery *v1.ZincQuery, query bluge.Query) bluge.SearchRequest {
	return bluge.NewTopNSearch(iQuery.MaxResults, query).
		SetFrom(iQuery.From).
		SortBy(iQuery.SortFields).
		WithStandardAggregations()
}
