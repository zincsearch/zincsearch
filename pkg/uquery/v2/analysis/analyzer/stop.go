package analyzer

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/analyzer"

	"github.com/prabhatsharma/zinc/pkg/bluge/analysis/token"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewStopAnalyzer(options interface{}) (*analysis.Analyzer, error) {
	stopwords, _ := zutils.GetStringSliceFromMap(options, "stopwords")
	ana := analyzer.NewSimpleAnalyzer()
	ana.TokenFilters = append(ana.TokenFilters, token.NewStopTokenFilter(stopwords))

	return ana, nil
}
