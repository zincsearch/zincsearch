package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func WildcardQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[wildcard] query doesn't support multiple fields")
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
				k := strings.ToLower(k)
				switch k {
				case "value":
					value.Value = v.(string)
				case "boost":
					value.Boost = v.(float64)
				default:
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[wildcard] unknown field [%s]", k))
				}
			}
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[wildcard] %s doesn't support values of type: %T", k, v))
		}
	}

	subq := bluge.NewWildcardQuery(value.Value).SetField(field)
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}
