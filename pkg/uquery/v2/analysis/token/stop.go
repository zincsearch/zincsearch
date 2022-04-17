package token

import (
	"github.com/blugelabs/bluge/analysis"

	"github.com/zinclabs/zinc/pkg/bluge/analysis/token"
	"github.com/zinclabs/zinc/pkg/zutils"
)

func NewStopTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	stopwords, _ := zutils.GetStringSliceFromMap(options, "stopwords")
	return token.NewStopTokenFilter(stopwords), nil
}
