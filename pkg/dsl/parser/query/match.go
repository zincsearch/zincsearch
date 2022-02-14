package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis/analyzer"

	"github.com/prabhatsharma/zinc/pkg/dsl/meta"
)

func MatchQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[match] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.MatchQuery)
	value.Boost = -1.0
	for k, v := range query {
		field = k
		switch v := v.(type) {
		case string:
			value.Query = v
		case map[string]interface{}:
			for k, v := range v {
				k := strings.ToLower(k)
				switch k {
				case "query":
					value.Query = v.(string)
				case "analyzer":
					value.Analyzer = v.(string)
				case "operator":
					value.Operator = v.(string)
				case "fuzziness":
					value.Fuzziness = v.(string)
				case "prefix_length":
					value.PrefixLength = v.(float64)
				case "boost":
					value.Boost = v.(float64)
				default:
					return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[match] unknown field [%s]", k))
				}
			}
		default:
			return nil, meta.NewError(meta.ErrorTypeXContentParseException, fmt.Sprintf("[match] %s doesn't support values of type: %T", k, v))
		}
	}

	subq := bluge.NewMatchQuery(value.Query).SetField(field)
	if value.Analyzer != "" {
		switch value.Analyzer {
		case "standard":
			subq.SetAnalyzer(analyzer.NewStandardAnalyzer())
		default:
			// TODO: support analyzer
		}
	}
	if value.Operator != "" {
		op := strings.ToUpper(value.Operator)
		switch op {
		case "OR":
			subq.SetOperator(bluge.MatchQueryOperatorOr)
		case "AND":
			subq.SetOperator(bluge.MatchQueryOperatorAnd)
		default:
			return nil, meta.NewError(meta.ErrorTypeIllegalArgumentException, fmt.Sprintf("[match] unknown operator %s", op))
		}
	}
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
