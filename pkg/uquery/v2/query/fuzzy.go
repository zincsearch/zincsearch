package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func FuzzyQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[fuzzy] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.FuzzyQuery)
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
				case "fuzziness":
					value.Fuzziness = v.(string)
				case "prefix_length":
					value.PrefixLength = v.(float64)
				case "boost":
					value.Boost = v.(float64)
				default:
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[fuzzy] unknown field [%s]", k))
				}
			}
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[fuzzy] %s doesn't support values of type: %T", k, v))
		}
	}

	subq := bluge.NewFuzzyQuery(value.Value).SetField(field)
	if value.Fuzziness != nil {
		switch v := value.Fuzziness.(type) {
		case string:
			// TODO: support other fuzziness: AUTO
		case float64:
			subq.SetFuzziness(int(v))
		}
	}
	if value.PrefixLength > 0 {
		subq.SetPrefix(int(value.PrefixLength))
	}
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}
