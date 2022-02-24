package analyzer

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/analyzer"
	"github.com/blugelabs/bluge/analysis/token"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewStandardAnalyzer(options interface{}) (*analysis.Analyzer, error) {
	stopwords, _ := zutils.GetStringSliceFromMap(options, "stopwords")
	maxTokenLength, _ := zutils.GetFloatFromMap(options, "max_token_length")
	if maxTokenLength == 0 {
		maxTokenLength = 255
	}

	ana := analyzer.NewStandardAnalyzer()
	if len(stopwords) > 0 {
		dict := analysis.NewTokenMap()
		for _, word := range stopwords {
			dict.AddToken(word)
		}
		ana.TokenFilters = append(ana.TokenFilters, token.NewStopTokensFilter(dict))
	}

	ana.TokenFilters = append(ana.TokenFilters, token.NewLengthFilter(1, int(maxTokenLength)))

	return ana, nil
}
