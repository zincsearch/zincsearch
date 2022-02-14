package query

import (
	"github.com/blugelabs/bluge"

	"github.com/prabhatsharma/zinc/pkg/dsl/meta"
)

func BoostingQuery(query map[string]interface{}) (bluge.Query, error) {
	return nil, meta.NewError(meta.ErrorTypeNotImplemented, "[boosting] query doesn't support")
}
