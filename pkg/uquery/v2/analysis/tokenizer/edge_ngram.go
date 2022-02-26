package tokenizer

import (
	"github.com/blugelabs/bluge/analysis"

	zinctokenizer "github.com/prabhatsharma/zinc/pkg/bluge/analysis/tokenizer"
	"github.com/prabhatsharma/zinc/pkg/errors"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewEdgeNgramTokenizer(options interface{}) (analysis.Tokenizer, error) {
	min, _ := zutils.GetFloatFromMap(options, "min_gram")
	max, _ := zutils.GetFloatFromMap(options, "max_gram")
	tokenChars, _ := zutils.GetStringSliceFromMap(options, "token_chars")
	if min == 0 {
		min = 1
	}
	if max == 0 {
		max = 2
	}
	if min > max {
		return nil, errors.New(errors.ErrorTypeParsingException, "[tokenizer] edge_ngram option [min_gram] should be not greater than [max_gram]")
	}
	return zinctokenizer.NewEdgeNgramTokenizer(int(min), int(max), tokenChars), nil
}
