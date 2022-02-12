package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func IdsQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[ids] query doesn't support multiple fields")
	}

	value := new(meta.IdsQuery)
	for k, v := range query {
		switch v := v.(type) {
		case []string:
			value.Values = v
		case []interface{}:
			value.Values = make([]string, len(v))
			for i, v := range v {
				value.Values[i] = v.(string)
			}
		case map[string]interface{}:
			for k, v := range v {
				k := strings.ToLower(k)
				switch k {
				case "value":
					switch v := v.(type) {
					case []interface{}:
						value.Values = make([]string, len(v))
						for i, v := range v {
							value.Values[i] = v.(string)
						}
					default:
						return nil, meta.NewError(meta.ErrorTypeXContentParseException, fmt.Sprintf("[ids] %s doesn't support values of type: %T", k, v))
					}
				default:
					return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[ids] unknown field [%s]", k))
				}
			}
		default:
			return nil, meta.NewError(meta.ErrorTypeXContentParseException, fmt.Sprintf("[ids] %s doesn't support values of type: %T", k, v))
		}
	}

	return TermsQuery(map[string]interface{}{
		"_id": value.Values,
	})
}
