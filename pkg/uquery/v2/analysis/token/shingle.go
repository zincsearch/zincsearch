package token

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewShingleTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	min, _ := zutils.GetFloatFromMap(options, "min_shingle_size")
	max, _ := zutils.GetFloatFromMap(options, "max_shingle_size")
	sep, _ := zutils.GetStringFromMap(options, "token_separator")
	fill, _ := zutils.GetStringFromMap(options, "filler_token")
	outputOriginalBool := true
	outputOriginal, _ := zutils.GetAnyFromMap(options, "output_original")
	if outputOriginal != nil {
		outputOriginalBool = outputOriginal.(bool)
	}
	if min == 0 {
		min = 2
	}
	if max == 0 {
		max = 2
	}
	if sep == "" {
		sep = " "
	}
	if fill == "" {
		fill = "_"
	}
	return token.NewShingleFilter(int(min), int(max), outputOriginalBool, sep, fill), nil
}
