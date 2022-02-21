package query

import (
	"github.com/blugelabs/bluge"
	"github.com/prabhatsharma/zinc/pkg/errors"
)

func BoostingQuery(query map[string]interface{}) (bluge.Query, error) {
	return nil, errors.New(errors.ErrorTypeNotImplemented, "[boosting] query doesn't support")
}
