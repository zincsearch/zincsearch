package v2

import (
	"github.com/blugelabs/bluge/search"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/parser/aggregation"
)

func FormatResponse(resp *meta.SearchResponse, q *meta.ZincQuery, buckets *search.Bucket) error {
	var err error
	// format aggregations
	if len(q.Aggregations) > 0 {
		resp.Aggregations, err = aggregation.Response(buckets)
		if err != nil {
			return meta.NewError(meta.ErrorTypeParsingException, err.Error())
		}
		if len(resp.Aggregations) > 0 {
			delete(resp.Aggregations, "count")
			delete(resp.Aggregations, "duration")
			delete(resp.Aggregations, "max_score")
		}
	}

	return nil
}
