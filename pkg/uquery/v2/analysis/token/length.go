package token

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewLengthTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	min, _ := zutils.GetFloatFromMap(options, "min")
	max, _ := zutils.GetFloatFromMap(options, "max")
	if min == 0 {
		min = 1
	}
	if max == 0 {
		max = 2
	}
	return token.NewLengthFilter(int(min), int(max)), nil
}
