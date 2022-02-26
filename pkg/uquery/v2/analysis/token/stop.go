package token

import (
	"github.com/blugelabs/bluge/analysis"

	"github.com/prabhatsharma/zinc/pkg/bluge/analysis/token"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewStopTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	stopwords, _ := zutils.GetStringSliceFromMap(options, "stopwords")
	return token.NewStopTokenFilter(stopwords), nil
}
