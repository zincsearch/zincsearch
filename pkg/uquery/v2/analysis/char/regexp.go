package char

import (
	"fmt"
	"regexp"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/char"

	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/zutils"
)

func NewRegexpCharFilter(options interface{}) (analysis.CharFilter, error) {
	pattern, err := zutils.GetStringFromMap(options, "pattern")
	if err != nil || pattern == "" {
		return nil, errors.New(errors.ErrorTypeParsingException, "[char_filter] regexp option [pattern] should be exists")
	}
	replacement, _ := zutils.GetStringFromMap(options, "replacement")
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[char_filter] regexp option [pattern] compile error: %s", err.Error()))
	}

	return char.NewRegexpCharFilter(r, []byte(replacement)), nil
}
