package analyzer

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/go-ego/gse"

	"github.com/prabhatsharma/zinc/pkg/bluge/analysis/lang/chs/token"
	"github.com/prabhatsharma/zinc/pkg/bluge/analysis/lang/chs/tokenizer"
)

func NewSearchAnalyzer(seg *gse.Segmenter) *analysis.Analyzer {
	return &analysis.Analyzer{
		Tokenizer:    tokenizer.NewSearchTokenizer(seg),
		TokenFilters: []analysis.TokenFilter{token.NewStopTokenFilter(seg, nil)},
	}
}
