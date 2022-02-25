package token

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/prabhatsharma/zinc/pkg/errors"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewKeywordTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	keywords, err := zutils.GetStringSliceFromMap(options, "keywords")
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, "[token_filter] keyword option [keywords] should be an array of string")
	}
	dict := analysis.NewTokenMap()
	for _, word := range keywords {
		dict.AddToken(word)
	}
	return token.NewKeyWordMarkerFilter(dict), nil
}
