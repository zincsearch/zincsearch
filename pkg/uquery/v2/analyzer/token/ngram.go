package token

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"

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
	return token.NewNgramFilter(int(min), int(max)), nil
}
