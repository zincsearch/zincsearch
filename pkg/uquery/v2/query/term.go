package query

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blugelabs/bluge"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func TermQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[term] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.TermQuery)
	value.Boost = -1.0
	for k, v := range query {
		field = k
		switch v := v.(type) {
		case string:
			value.Value = v
		case float64:
			value.Value = v
		case bool:
			value.Value = v
		case map[string]interface{}:
			for k, v := range v {
				k := strings.ToLower(k)
				switch k {
				case "value":
					switch vv := v.(type) {
					case string:
						value.Value = vv
					case float64:
						value.Value = vv
					case bool:
						value.Value = vv
					default:
						return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[term] doesn't support values of type: %T", v))
					}
				case "case_insensitive":
					value.CaseInsensitive = v.(bool)
				case "boost":
					value.Boost = v.(float64)
				default:
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[term] unknown field [%s]", k))
				}
			}
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[term] doesn't support values of type: %T", v))
		}
	}

	// TODO: case_insensitive support

	switch value.Value.(type) {
	case string:
		subq := bluge.NewTermQuery(value.Value.(string)).SetField(field)
		if value.Boost >= 0 {
			subq.SetBoost(value.Boost)
		}
		return subq, nil
	case float64:
		subq := bluge.NewNumericRangeInclusiveQuery(value.Value.(float64), value.Value.(float64), true, true).SetField(field)
		if value.Boost >= 0 {
			subq.SetBoost(value.Boost)
		}
		return subq, nil
	case bool:
		subq := bluge.NewTermQuery(strconv.FormatBool(value.Value.(bool))).SetField(field)
		if value.Boost >= 0 {
			subq.SetBoost(value.Boost)
		}
		return subq, nil
	default:
		return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[term] doesn't support values of type: %T", value.Value))
	}
}
