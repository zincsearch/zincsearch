package query

import (
	"github.com/blugelabs/bluge"

	"github.com/prabhatsharma/zinc/pkg/dsl/meta"
)

func CombinedFieldsQuery(query map[string]interface{}) (bluge.Query, error) {
	return nil, meta.NewError(meta.ErrorTypeNotImplemented, "[combined_fields] query doesn't support")
}
