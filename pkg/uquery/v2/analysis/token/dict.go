package token

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/prabhatsharma/zinc/pkg/errors"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewDictTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	words, err := zutils.GetStringSliceFromMap(options, "words")
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, "[token_filter] dict option [words] should be an array of string")
	}
	dict := analysis.NewTokenMap()
	for _, word := range words {
		dict.AddToken(word)
	}

	minWordSize, _ := zutils.GetFloatFromMap(options, "min_word_size")
	minSubWordSize, _ := zutils.GetFloatFromMap(options, "min_sub_word_size")
	maxSubWordSize, _ := zutils.GetFloatFromMap(options, "max_sub_word_size")
	onlyLongestMatch, _ := zutils.GetBoolFromMap(options, "only_longest_match")
	if minWordSize == 0 {
		minWordSize = 5
	}
	if minSubWordSize == 0 {
		minSubWordSize = 2
	}
	if maxSubWordSize == 0 {
		maxSubWordSize = 15
	}
	return token.NewDictionaryCompoundFilter(dict, int(minWordSize), int(minSubWordSize), int(maxSubWordSize), onlyLongestMatch), nil
}
