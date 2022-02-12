package query

import (
	"github.com/blugelabs/bluge"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func ExistsQuery(query map[string]interface{}) (bluge.Query, error) {
	return nil, meta.NewError(meta.ErrorTypeNotImplemented, "[exists] query doesn't support")
}
