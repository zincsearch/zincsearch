package query

import (
	"github.com/blugelabs/bluge"
	"github.com/prabhatsharma/zinc/pkg/errors"
)

func CombinedFieldsQuery(query map[string]interface{}) (bluge.Query, error) {
	return nil, errors.New(errors.ErrorTypeNotImplemented, "[combined_fields] query doesn't support")
}
