package query

import (
	"github.com/blugelabs/bluge"
	"github.com/prabhatsharma/zinc/pkg/errors"
)

func TermsSetQuery(query map[string]interface{}) (bluge.Query, error) {
	return nil, errors.New(errors.ErrorTypeNotImplemented, "[terms_set] query doesn't support")
}
