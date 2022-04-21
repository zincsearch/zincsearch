package tokenizer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/tokenizer"

	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/zutils"
)

func NewExceptionTokenizer(options interface{}) (analysis.Tokenizer, error) {
	patterns, err := zutils.GetStringSliceFromMap(options, "patterns")
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, "[tokenizer] exception option [patterns] should be an array of string")
	}

	pattern := strings.Join(patterns, "|")
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[tokenizer] exception option [patterns] compile error: %s", err.Error()))
	}

	return tokenizer.NewExceptionsTokenizer(r, tokenizer.NewUnicodeTokenizer()), nil
}
