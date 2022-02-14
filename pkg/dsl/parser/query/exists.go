package query

import (
	"github.com/blugelabs/bluge"

	"github.com/prabhatsharma/zinc/pkg/dsl/meta"
)

func ExistsQuery(query map[string]interface{}) (bluge.Query, error) {
	return nil, meta.NewError(meta.ErrorTypeNotImplemented, "[exists] query doesn't support")
}
