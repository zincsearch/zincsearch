package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	zincanalysis "github.com/prabhatsharma/zinc/pkg/uquery/v2/analysis"
)

func MatchPhraseQuery(query map[string]interface{}, mappings *meta.Mappings, analyzers map[string]*analysis.Analyzer) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[match_phrase] query doesn't support multiple fields")
	}

	field := ""
	value := new(meta.MatchPhraseQuery)
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
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[match_phrase] unknown field [%s]", k))
				}
			}
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[match_phrase] %s doesn't support values of type: %T", k, v))
		}
	}

	var err error
	var zer *analysis.Analyzer
	if value.Analyzer != "" {
		zer, err = zincanalysis.QueryAnalyzer(analyzers, value.Analyzer)
		if err != nil {
			return nil, err
		}
	} else {
		indexZer, searchZer := zincanalysis.QueryAnalyzerForField(analyzers, mappings, field)
		if zer == nil && searchZer != nil {
			zer = searchZer
		}
		if zer == nil && indexZer != nil {
			zer = indexZer
		}
	}

	subq := bluge.NewMatchPhraseQuery(value.Query).SetField(field)
	if zer != nil {
		subq.SetAnalyzer(zer)
	}
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}
