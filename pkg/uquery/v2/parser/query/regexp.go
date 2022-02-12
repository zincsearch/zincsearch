package query

import (
	"github.com/blugelabs/bluge"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func RegexpQuery(query map[string]interface{}) (bluge.Query, error) {
	return nil, meta.NewError(meta.ErrorTypeNotImplemented, "[regexp] query doesn't support")
}
