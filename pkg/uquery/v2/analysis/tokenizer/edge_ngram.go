package tokenizer

import (
	"github.com/blugelabs/bluge/analysis"

	zinctokenizer "github.com/prabhatsharma/zinc/pkg/bluge/analysis/tokenizer"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewEdgeNgramTokenizer(options interface{}) (analysis.Tokenizer, error) {
	min, _ := zutils.GetFloatFromMap(options, "min_gram")
	max, _ := zutils.GetFloatFromMap(options, "max_gram")
	if min == 0 {
		min = 1
	}
	if max == 0 {
		max = 2
	}
	return zinctokenizer.NewEdgeNgramTokenizer(int(min), int(max)), nil
}
