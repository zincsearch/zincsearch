package analyzer

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/analyzer"
	"github.com/blugelabs/bluge/analysis/char"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func Request(data *meta.IndexAnalysis) (map[string]*analysis.Analyzer, error) {
	if data == nil {
		return nil, nil
	}

	if data.Analyzer == nil {
		return nil, nil
	}

	charFilters, err := RequestCharFilter(data.CharFilter)
	if err != nil {
		return nil, err
	}

	tokenFilters, err := RequestTokenFilter(data.TokenFilter)
	if err != nil {
		return nil, err
	}

	analyzers := make(map[string]*analysis.Analyzer)
	for name, v := range data.Analyzer {
		var ana *analysis.Analyzer
		if v.Tokenizer == "" {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] is missing tokenizer", name))
		}

		chars := make([]analysis.CharFilter, 0, len(v.CharFilter))
		for _, filter := range v.CharFilter {
			switch filter {
			case "ascii_folding":
				chars = append(chars, char.NewASCIIFoldingFilter())
			case "html", "html_strip":
				chars = append(chars, char.NewHTMLCharFilter())
			case "zero_width_non_joiner":
				chars = append(chars, char.NewZeroWidthNonJoinerCharFilter())
			default:
				if v, ok := charFilters[filter]; ok {
					chars = append(chars, v)
				} else {
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] used undefined char_filter [%s]", name, filter))
				}
			}
		}

		tokens := make([]analysis.TokenFilter, 0, len(v.TokenFilter))
		for _, filter := range v.TokenFilter {
			switch filter {
			case "apostrophe":
				tokens = append(tokens, token.NewApostropheFilter())
			case "camel_case":
				tokens = append(tokens, token.NewCamelCaseFilter())
			case "lower_case":
				tokens = append(tokens, token.NewLowerCaseFilter())
			case "porter":
				tokens = append(tokens, token.NewPorterStemmer())
			case "reverse":
				tokens = append(tokens, token.NewReverseFilter())
			case "unique":
				tokens = append(tokens, token.NewUniqueTermFilter())
			default:
				if v, ok := tokenFilters[filter]; ok {
					tokens = append(tokens, v)
				} else {
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] used undefined token_filter [%s]", name, filter))
				}
			}
		}

		v.Tokenizer = strings.ToLower(v.Tokenizer)
		switch v.Tokenizer {
		case "standard":
			ana = analyzer.NewStandardAnalyzer()
		case "keyword":
			ana = analyzer.NewKeywordAnalyzer()
		case "simple":
			ana = analyzer.NewSimpleAnalyzer()
		case "web":
			ana = analyzer.NewWebAnalyzer()
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[analyzer] [%s] doesn't support tokenizer [%s]", name, v.Tokenizer))
		}

		if len(chars) > 0 {
			ana.CharFilters = append(ana.CharFilters, chars...)
		}
		if len(tokens) > 0 {
			ana.TokenFilters = append(ana.TokenFilters, tokens...)
		}
		analyzers[name] = ana
	}

	return analyzers, nil
}

func Query(data map[string]*analysis.Analyzer, name string) (*analysis.Analyzer, error) {
	if data != nil {
		if v, ok := data[name]; ok {
			return v, nil
		}
	}

	switch name {
	case "", "standard":
		return analyzer.NewStandardAnalyzer(), nil
	case "keyword":
		return analyzer.NewKeywordAnalyzer(), nil
	case "simple":
		return analyzer.NewSimpleAnalyzer(), nil
	case "web":
		return analyzer.NewWebAnalyzer(), nil
	default:
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] doesn't exists", name))
	}
}
