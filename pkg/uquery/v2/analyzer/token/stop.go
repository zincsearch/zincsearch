package token

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/prabhatsharma/zinc/pkg/errors"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewStopTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	stopwords, err := zutils.GetStringSliceFromMap(options, "stopwords")
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, "[token_filter] keyword option [keywords] should be an array of strings")
	}
	dict := analysis.NewTokenMap()
	for _, word := range stopwords {
		dict.AddToken(word)
	}
	return token.NewStopTokensFilter(dict), nil
}
