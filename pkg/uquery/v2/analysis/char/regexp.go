package char

import (
	"fmt"
	"regexp"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/char"

	"github.com/prabhatsharma/zinc/pkg/errors"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewRegexpCharFilter(options interface{}) (analysis.CharFilter, error) {
	pattern, err := zutils.GetStringFromMap(options, "pattern")
	if err != nil || pattern == "" {
		return nil, errors.New(errors.ErrorTypeParsingException, "[char_filter] regexp option [pattern] should be exists")
	}
	replacement, err := zutils.GetStringFromMap(options, "replacement")
	if err != nil || replacement == "" {
		return nil, errors.New(errors.ErrorTypeParsingException, "[char_filter] regexp option [replacement] should be exists")
	}
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[char_filter] regexp option [pattern] compile error: %v", err.Error()))
	}

	return char.NewRegexpCharFilter(r, []byte(replacement)), nil
}
