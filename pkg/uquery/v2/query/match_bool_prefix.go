package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis/analyzer"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func MatchBoolPrefixQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[match_bool_prefix] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.MatchBoolPrefixQuery)
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
				default:
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[match_bool_prefix] unknown field [%s]", k))
				}
			}
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[match_bool_prefix] %s doesn't support values of type: %T", k, v))
		}
	}

	// TODO support analyzer
	zer := analyzer.NewStandardAnalyzer()
	if value.Analyzer != "" {
		switch value.Analyzer {
		case "standard":
			zer = analyzer.NewStandardAnalyzer()
		default:
			// TODO: support analyzer
		}
	}

	tokens := zer.Analyze([]byte(value.Query))
	subq := bluge.NewBooleanQuery()
	for i := 0; i < len(tokens); i++ {
		if i == len(tokens)-1 {
			subq.AddShould(bluge.NewPrefixQuery(string(tokens[i].Term)).SetField(field))
		} else {
			subq.AddShould(bluge.NewTermQuery(string(tokens[i].Term)).SetField(field))
		}
	}

	return subq, nil
}
