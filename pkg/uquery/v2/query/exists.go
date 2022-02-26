package query

import (
	"github.com/blugelabs/bluge"
	"github.com/prabhatsharma/zinc/pkg/errors"
)

func ExistsQuery(query map[string]interface{}) (bluge.Query, error) {
	return nil, errors.New(errors.ErrorTypeNotImplemented, "[exists] query doesn't support")
}
