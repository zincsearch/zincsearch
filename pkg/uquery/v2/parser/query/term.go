package query

import (
	"github.com/blugelabs/bluge"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func TermQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[term] query doesn't support multiple fields")
	}

	return nil, nil
}
