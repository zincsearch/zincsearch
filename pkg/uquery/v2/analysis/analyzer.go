package analysis

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/analyzer"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	zincanalyzer "github.com/prabhatsharma/zinc/pkg/uquery/v2/analysis/analyzer"
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

	tokenizers, err := RequestTokenizer(data.Tokenizer)
	if err != nil {
		return nil, err
	}

	analyzers := make(map[string]*analysis.Analyzer)
	for name, v := range data.Analyzer {
		if v.Tokenizer == "" {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] is missing tokenizer", name))
		}

		// custom build-in analyzer
		var ana *analysis.Analyzer
		if v.Type != "" {
			v.Type = strings.ToLower(v.Type)
			switch v.Type {
			case "pattern":
				ana, err = zincanalyzer.NewPatternAnalyzer(map[string]interface{}{
					"pattern":   v.Pattern,
					"lowercase": v.Lowercase,
					"stopwords": v.Stopwords,
				})
			case "standard":
				ana, err = zincanalyzer.NewStandardAnalyzer(map[string]interface{}{
					"max_token_length": v.MaxTokenLength,
					"stopwords":        v.Stopwords,
				})
			case "stop":
				ana, err = zincanalyzer.NewStopAnalyzer(map[string]interface{}{
					"stopwords": v.Stopwords,
				})
			default:
				return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] build-in [%s] doesn't support custom", v.Type))
			}
			if err != nil {
				return nil, err
			}
		}

		// use tokenizer
		var ok bool
		zer, err := RequestTokenizerSingle(v.Tokenizer, nil)
		if zer != nil && err == nil {
			// use standard tokenizer
		} else {
			if zer, ok = tokenizers[v.Tokenizer]; !ok {
				if ana == nil { // returns error if not user build-in analyzer
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] used undifined tokenizer %s", name, v.Tokenizer))
				}
			}
		}

		chars := make([]analysis.CharFilter, 0, len(v.CharFilter))
		for _, name := range v.CharFilter {
			filter, err := RequestCharFilterSingle(name, nil)
			if filter != nil && err == nil {
				chars = append(chars, filter)
			} else {
				if v, ok := charFilters[name]; ok {
					chars = append(chars, v)
				} else {
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] used undefined char_filter [%s]", name, filter))
				}
			}
		}

		tokens := make([]analysis.TokenFilter, 0, len(v.TokenFilter))
		for _, name := range v.TokenFilter {
			filter, err := RequestTokenFilterSingle(name, nil)
			if filter != nil && err == nil {
				tokens = append(tokens, filter)
			} else {
				if v, ok := tokenFilters[name]; ok {
					tokens = append(tokens, v)
				} else {
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] used undefined token_filter [%s]", name, filter))
				}
			}
		}

		if ana == nil {
			ana = &analysis.Analyzer{Tokenizer: zer}
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
		return zincanalyzer.NewStandardAnalyzer(nil)
	case "keyword":
		return analyzer.NewKeywordAnalyzer(), nil
	case "simple":
		return analyzer.NewSimpleAnalyzer(), nil
	case "web":
		return analyzer.NewWebAnalyzer(), nil
	case "pattern":
		return zincanalyzer.NewPatternAnalyzer(nil)
	case "whitespace":
		return zincanalyzer.NewWhitespaceAnalyzer()
	case "stop":
		return zincanalyzer.NewStopAnalyzer(nil)
	default:
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] doesn't exists", name))
	}
}
