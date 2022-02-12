package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis/analyzer"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func MatchPhrasePrefixQuery(query map[string]interface{}) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "[match_phrase_prefix] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.MatchPhrasePrefixQuery)
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
				case "boost":
					value.Boost = v.(float64)
				default:
					return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[match_phrase_prefix] unknown field [%s]", k))
				}
			}
		default:
			return nil, meta.NewError(meta.ErrorTypeXContentParseException, fmt.Sprintf("[match_phrase_prefix] %s doesn't support values of type: %T", k, v))
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
	if len(tokens) > 0 {
		subq.AddMust(bluge.NewPrefixQuery(string(tokens[len(tokens)-1].Term)).SetField(field))
	}
	if len(tokens) > 1 {
		phrase := strings.Builder{}
		for i := 0; i < len(tokens)-1; i++ {
			phrase.WriteString(string(tokens[i].Term))
			phrase.WriteString(" ")
		}
		subq.AddMust(bluge.NewMatchPhraseQuery(strings.TrimSpace(phrase.String())).SetField(field).SetAnalyzer(zer))
	}
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}
