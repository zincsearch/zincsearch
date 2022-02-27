package analyzer

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/tokenizer"
)

func NewWhitespaceAnalyzer() (*analysis.Analyzer, error) {
	return &analysis.Analyzer{Tokenizer: tokenizer.NewWhitespaceTokenizer()}, nil
}
