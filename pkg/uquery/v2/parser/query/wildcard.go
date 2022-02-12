package query

import (
	"fmt"

	"github.com/blugelabs/bluge"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func WildcardQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[wildcard] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.WildcardQuery)
	value.Boost = -1.0
	for k, v := range query {
		field = k
		switch v := v.(type) {
		case string:
			value.Value = v
		case map[string]interface{}:
			for k, v := range v {
				switch k {
				case "value":
					value.Value = v.(string)
				case "boost":
					value.Boost = v.(float64)
				default:
					return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[wildcard] unsupported children %s", k))
				}
			}
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[wildcard] unsupported query type %s", k))
		}
	}

	subq := bluge.NewWildcardQuery(value.Value).SetField(field)
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}
