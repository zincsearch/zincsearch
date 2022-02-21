package analyzer

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/analyzer"
	"github.com/blugelabs/bluge/analysis/char"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func Request(data *meta.IndexAnalysis) (map[string]*analysis.Analyzer, error) {
	if data == nil {
		return nil, nil
	}

	charFilters, err := QueryCharFilter(data.CharFilter)
	if err != nil {
		return nil, err
	}

	if data.Analyzer == nil {
		return nil, nil
	}

	analyzers := make(map[string]*analysis.Analyzer)
	for name, v := range data.Analyzer {
		var ana *analysis.Analyzer
		if v.Tokenizer == "" {
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] is missing tokenizer", name))
		}

		chars := make([]analysis.CharFilter, 0, len(v.CharFilter))
		for _, filter := range v.CharFilter {
			switch filter {
			case "ascii_folding":
				chars = append(chars, char.NewASCIIFoldingFilter())
			case "html":
				chars = append(chars, char.NewHTMLCharFilter())
			case "zero_width_non_joiner":
				chars = append(chars, char.NewZeroWidthNonJoinerCharFilter())
			default:
				if v, ok := charFilters[filter]; ok {
					chars = append(chars, v)
				} else {
					return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] used undefined char_filter [%s]", name, filter))
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
			return nil, meta.NewError(meta.ErrorTypeXContentParseException, fmt.Sprintf("[analyzer] [%s] doesn't support tokenizer [%s]", name, v.Tokenizer))
		}

		if len(chars) > 0 {
			ana.CharFilters = append(ana.CharFilters, chars...)
		}
		analyzers[name] = ana
	}

	return analyzers, nil
}

func Query(data map[string]*analysis.Analyzer, name string) (*analysis.Analyzer, error) {
	if v, ok := data[name]; ok {
		return v, nil
	}

	switch name {
	case "standard":
		return analyzer.NewStandardAnalyzer(), nil
	case "keyword":
		return analyzer.NewKeywordAnalyzer(), nil
	case "simple":
		return analyzer.NewSimpleAnalyzer(), nil
	case "web":
		return analyzer.NewWebAnalyzer(), nil
	default:
		return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] doesn't exists", name))
	}
}
