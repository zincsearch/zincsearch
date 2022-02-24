package token

import (
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewEdgeNgramTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	min, _ := zutils.GetFloatFromMap(options, "min_gram")
	max, _ := zutils.GetFloatFromMap(options, "max_gram")
	side, _ := zutils.GetStringFromMap(options, "side")
	boolSide := token.Side(true)
	side = strings.ToLower(side)
	if side != "back" {
		side = "front"
		boolSide = false
	}
	if min == 0 {
		min = 1
	}
	if max == 0 {
		max = 2
	}
	return token.NewEdgeNgramFilter(boolSide, int(min), int(max)), nil
}
