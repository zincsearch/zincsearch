package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func PrefixQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[prefix] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.PrefixQuery)
	value.Boost = -1.0
	for k, v := range query {
		field = k
		switch v := v.(type) {
		case string:
			value.Value = v
		case map[string]interface{}:
			for k, v := range v {
				k := strings.ToLower(k)
				switch k {
				case "value":
					value.Value = v.(string)
				case "boost":
					value.Boost = v.(float64)
				default:
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[prefix] unknown field [%s]", k))
				}
			}
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[prefix] %s doesn't support values of type: %T", k, v))
		}
	}

	subq := bluge.NewPrefixQuery(value.Value).SetField(field)
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}
