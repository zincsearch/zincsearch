package token

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/prabhatsharma/zinc/pkg/errors"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewElisionTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	articles, err := zutils.GetStringSliceFromMap(options, "articles")
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, "[token_filter] keyword option [keywords] should be an array of strings")
	}
	dict := analysis.NewTokenMap()
	for _, keyword := range articles {
		dict.AddToken(keyword)
	}
	return token.NewElisionFilter(dict), nil
}
