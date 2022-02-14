package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"

	"github.com/prabhatsharma/zinc/pkg/dsl/meta"
)

func TermsQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 2 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[terms] query doesn't support multiple fields")
	}

	field := ""
	values := []string{}
	boost := -1.0
	for k, v := range query {
		if strings.ToLower(k) == "boost" {
			boost = v.(float64)
			continue
		}

		field = k
		switch v := v.(type) {
		case []string:
			values = v
		case []interface{}:
			for _, vv := range v {
				values = append(values, vv.(string))
			}
		default:
			return nil, meta.NewError(meta.ErrorTypeXContentParseException, fmt.Sprintf("[terms] %s doesn't support values of type: %T", k, v))
		}
	}

	subq := bluge.NewBooleanQuery()
	for _, term := range values {
		subq.AddShould(bluge.NewTermQuery(term).SetField(field))
	}
	if boost >= 0 {
		subq.SetBoost(boost)
	}

	return subq, nil
}
