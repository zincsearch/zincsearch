package token

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"
)

func NewStopTokenFilter(stopwords []string) analysis.TokenFilter {
	rv := StopWords(stopwords)
	return token.NewStopTokensFilter(rv)
}
