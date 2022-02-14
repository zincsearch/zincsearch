package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"

	"github.com/prabhatsharma/zinc/pkg/dsl/meta"
)

func RangeQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[range] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.RangeQuery)
	value.Boost = -1.0
	for k, v := range query {
		field = k
		switch v := v.(type) {
		case map[string]interface{}:
			for k, v := range v {
				k := strings.ToLower(k)
				switch k {
				case "gt":
					value.GT = v.(float64)
				case "gte":
					value.GTE = v.(float64)
				case "lt":
					value.LT = v.(float64)
				case "lte":
					value.LTE = v.(float64)
				case "format":
					value.Format = v.(string)
				case "time_zone":
					value.TimeZone = v.(string)
				case "boost":
					value.Boost = v.(float64)
				default:
					return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[range] unknown field [%s]", k))
				}
			}
		default:
			return nil, meta.NewError(meta.ErrorTypeXContentParseException, fmt.Sprintf("[range] %s doesn't support values of type: %T", k, v))
		}
	}

	// TODO: choose range type by field mappings

	min := 0.0
	max := 0.0
	minInclusive := false
	maxInclusive := false
	if value.GT != nil && value.GT.(float64) > 0 {
		min = value.GT.(float64)

	}
	if value.GTE != nil && value.GTE.(float64) > 0 {
		min = value.GTE.(float64)
		minInclusive = true
	}
	if value.LT != nil && value.LT.(float64) > 0 {
		max = value.LT.(float64)
	}
	if value.LTE != nil && value.LTE.(float64) > 0 {
		max = value.LTE.(float64)
		maxInclusive = true
	}
	subq := bluge.NewNumericRangeInclusiveQuery(min, max, minInclusive, maxInclusive).SetField(field)
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}
