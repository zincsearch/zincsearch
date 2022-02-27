package analyzer

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/analyzer"

	zinctoken "github.com/prabhatsharma/zinc/pkg/bluge/analysis/token"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewStandardAnalyzer(options interface{}) (*analysis.Analyzer, error) {
	stopwords, _ := zutils.GetStringSliceFromMap(options, "stopwords")

	ana := analyzer.NewStandardAnalyzer()
	if len(stopwords) > 0 {
		ana.TokenFilters = append(ana.TokenFilters, zinctoken.NewStopTokenFilter(stopwords))
	}

	return ana, nil
}
