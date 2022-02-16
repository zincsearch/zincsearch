package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func TermsQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 2 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[terms] query doesn't support multiple fields")
	}

	field := ""
	values := []string{}
	valueInts := []float64{}
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
		case []float64:
			valueInts = v
		case []interface{}:
			for _, vv := range v {
				switch vvv := vv.(type) {
				case string:
					values = append(values, vvv)
				case float64:
					valueInts = append(valueInts, vvv)
				default:
					return nil, meta.NewError(meta.ErrorTypeXContentParseException, fmt.Sprintf("[term] doesn't support values of type: %T", vv))
				}
			}
		default:
			return nil, meta.NewError(meta.ErrorTypeXContentParseException, fmt.Sprintf("[terms] doesn't support values of type: %T", v))
		}
	}

	subq := bluge.NewBooleanQuery()
	for _, term := range values {
		subq.AddShould(bluge.NewTermQuery(term).SetField(field))
	}
	for _, term := range valueInts {
		subq.AddShould(bluge.NewNumericRangeInclusiveQuery(term, term, true, true).SetField(field))
	}
	if boost >= 0 {
		subq.SetBoost(boost)
	}

	return subq, nil
}
