package query

import (
	"github.com/blugelabs/bluge"

	"github.com/prabhatsharma/zinc/pkg/dsl/meta"
)

func TermsSetQuery(query map[string]interface{}) (bluge.Query, error) {
	return nil, meta.NewError(meta.ErrorTypeNotImplemented, "[terms_set] query doesn't support")
}
