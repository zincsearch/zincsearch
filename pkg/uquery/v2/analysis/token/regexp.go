package token

import (
	"fmt"
	"regexp"

	"github.com/blugelabs/bluge/analysis"

	"github.com/prabhatsharma/zinc/pkg/bluge/analysis/token"
	"github.com/prabhatsharma/zinc/pkg/errors"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewRegexpTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	pattern, err := zutils.GetStringFromMap(options, "pattern")
	if err != nil || pattern == "" {
		return nil, errors.New(errors.ErrorTypeParsingException, "[token_filter] regexp option [pattern] should be exists")
	}
	replacement, _ := zutils.GetStringFromMap(options, "replacement")
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[token_filter] regexp option [pattern] compile error: %v", err.Error()))
	}

	return token.NewRegexpTokenFilter(r, []byte(replacement)), nil
}
