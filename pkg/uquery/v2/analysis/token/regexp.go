package token

import (
	"fmt"
	"regexp"

	"github.com/blugelabs/bluge/analysis"

	"github.com/zinclabs/zinc/pkg/bluge/analysis/token"
	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/zutils"
)

func NewRegexpTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	pattern, err := zutils.GetStringFromMap(options, "pattern")
	if err != nil || pattern == "" {
		return nil, errors.New(errors.ErrorTypeParsingException, "[token_filter] regexp option [pattern] should be exists")
	}
	replacement, _ := zutils.GetStringFromMap(options, "replacement")
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[token_filter] regexp option [pattern] compile error: %s", err.Error()))
	}

	return token.NewRegexpTokenFilter(r, []byte(replacement)), nil
}
