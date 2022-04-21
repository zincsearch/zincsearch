package tokenizer

import (
	"fmt"
	"regexp"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/tokenizer"
	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/zutils"
)

func NewRegexpTokenizer(options interface{}) (analysis.Tokenizer, error) {
	pattern, _ := zutils.GetStringFromMap(options, "pattern")
	if len(pattern) == 0 {
		pattern = "\\w+"
	}

	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[tokenizer] regexp option [pattern] compile error: %s", err.Error()))
	}

	return tokenizer.NewRegexpTokenizer(r), nil
}
