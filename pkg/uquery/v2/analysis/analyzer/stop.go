package analyzer

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/analyzer"

	"github.com/zinclabs/zinc/pkg/bluge/analysis/token"
	"github.com/zinclabs/zinc/pkg/zutils"
)

func NewStopAnalyzer(options interface{}) (*analysis.Analyzer, error) {
	stopwords, _ := zutils.GetStringSliceFromMap(options, "stopwords")
	ana := analyzer.NewSimpleAnalyzer()
	ana.TokenFilters = append(ana.TokenFilters, token.NewStopTokenFilter(stopwords))

	return ana, nil
}
