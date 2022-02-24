package analyzer

import (
	"fmt"
	"regexp"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"
	"github.com/blugelabs/bluge/analysis/tokenizer"
	"github.com/prabhatsharma/zinc/pkg/errors"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewPatternAnalyzer(options interface{}) (*analysis.Analyzer, error) {
	pattern, _ := zutils.GetStringFromMap(options, "pattern")
	if pattern == "" {
		pattern = "\\W+"
	}
	lowerCase, err := zutils.GetBoolFromMap(options, "lowercase")
	if err != nil {
		lowerCase = true
	}
	stopwords, _ := zutils.GetStringSliceFromMap(options, "stopwords")
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] pattern option [pattern] compile error: %v", err.Error()))
	}
	ana := &analysis.Analyzer{Tokenizer: tokenizer.NewRegexpTokenizer(r)}
	if lowerCase {
		ana.TokenFilters = append(ana.TokenFilters, token.NewLowerCaseFilter())
	}

	if len(stopwords) > 0 {
		dict := analysis.NewTokenMap()
		for _, word := range stopwords {
			dict.AddToken(word)
		}
		ana.TokenFilters = append(ana.TokenFilters, token.NewStopTokensFilter(dict))
	}

	return ana, nil
}
