package query

import (
	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func SimpleQueryStringQuery(query map[string]interface{}, mappings *meta.Mappings, analyzers map[string]*analysis.Analyzer) (bluge.Query, error) {
	return QueryStringQuery(query, mappings, analyzers)
}
