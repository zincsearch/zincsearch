package token

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/prabhatsharma/zinc/pkg/errors"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewNgramTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	min, _ := zutils.GetFloatFromMap(options, "min_gram")
	max, _ := zutils.GetFloatFromMap(options, "max_gram")
	if min == 0 {
		min = 1
	}
	if max == 0 {
		max = 2
	}
	if min > max {
		return nil, errors.New(errors.ErrorTypeParsingException, "[token_filter] ngram option [min_gram] should be not greater than [max_gram]")
	}
	return token.NewNgramFilter(int(min), int(max)), nil
}
