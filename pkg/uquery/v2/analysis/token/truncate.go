package token

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewTruncateTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	length, _ := zutils.GetFloatFromMap(options, "length")
	if length == 0 {
		length = 10
	}
	return token.NewTruncateTokenFilter(int(length)), nil
}
