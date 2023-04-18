/* Copyright 2022 Zinc Labs Inc. and Contributors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package analysis

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/analyzer"
	"github.com/blugelabs/bluge/analysis/lang/ar"
	"github.com/blugelabs/bluge/analysis/lang/cjk"
	"github.com/blugelabs/bluge/analysis/lang/ckb"
	"github.com/blugelabs/bluge/analysis/lang/da"
	"github.com/blugelabs/bluge/analysis/lang/de"
	"github.com/blugelabs/bluge/analysis/lang/en"
	"github.com/blugelabs/bluge/analysis/lang/es"
	"github.com/blugelabs/bluge/analysis/lang/fa"
	"github.com/blugelabs/bluge/analysis/lang/fi"
	"github.com/blugelabs/bluge/analysis/lang/fr"
	"github.com/blugelabs/bluge/analysis/lang/hi"
	"github.com/blugelabs/bluge/analysis/lang/hu"
	"github.com/blugelabs/bluge/analysis/lang/it"
	"github.com/blugelabs/bluge/analysis/lang/nl"
	"github.com/blugelabs/bluge/analysis/lang/no"
	"github.com/blugelabs/bluge/analysis/lang/pt"
	"github.com/blugelabs/bluge/analysis/lang/ro"
	"github.com/blugelabs/bluge/analysis/lang/ru"
	"github.com/blugelabs/bluge/analysis/lang/sv"
	"github.com/blugelabs/bluge/analysis/lang/tr"

	"github.com/zincsearch/zincsearch/pkg/bluge/analysis/lang/chs"
	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/meta"
	zincanalyzer "github.com/zincsearch/zincsearch/pkg/uquery/analysis/analyzer"
)

func RequestAnalyzer(data *meta.IndexAnalysis) (map[string]*analysis.Analyzer, error) {
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

	if data.TokenFilter == nil && data.Filter != nil {
		data.TokenFilter = data.Filter
		data.Filter = nil
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
		if v.Tokenizer == "" && v.Type == "" {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] is missing tokenizer", name))
		}

		// custom build-in analyzer
		var ana *analysis.Analyzer
		if v.Type != "" {
			v.Type = strings.ToLower(v.Type)
			switch v.Type {
			case "custom":
				// omit
			case "regexp", "pattern":
				ana, err = zincanalyzer.NewRegexpAnalyzer(map[string]interface{}{
					"pattern":   v.Pattern,
					"lowercase": v.Lowercase,
					"stopwords": v.Stopwords,
				})
			case "standard":
				ana, err = zincanalyzer.NewStandardAnalyzer(map[string]interface{}{
					"stopwords": v.Stopwords,
				})
			case "stop":
				ana, err = zincanalyzer.NewStopAnalyzer(map[string]interface{}{
					"stopwords": v.Stopwords,
				})
			default:
				ana, err = QueryAnalyzer(nil, v.Type)
				if ana == nil {
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] unsuported build-in analyzer [%s]", v.Type))
				}
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
		for _, filterName := range v.CharFilter {
			filter, err := RequestCharFilterSingle(filterName, nil)
			if filter != nil && err == nil {
				chars = append(chars, filter)
			} else {
				if v, ok := charFilters[filterName]; ok {
					chars = append(chars, v)
				} else {
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] used undefined char_filter [%s]", name, filterName))
				}
			}
		}

		tokens := make([]analysis.TokenFilter, 0, len(v.TokenFilter))
		if v.TokenFilter == nil && v.Filter != nil {
			v.TokenFilter = v.Filter
			v.Filter = nil
		}
		for _, filterName := range v.TokenFilter {
			filter, err := RequestTokenFilterSingle(filterName, nil)
			if filter != nil && err == nil {
				tokens = append(tokens, filter)
			} else {
				if v, ok := tokenFilters[filterName]; ok {
					tokens = append(tokens, v)
				} else {
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] [%s] used undefined token_filter [%s]", name, filterName))
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

func QueryAnalyzer(data map[string]*analysis.Analyzer, name string) (*analysis.Analyzer, error) {
	if name == "" {
		name = "default"
	}

	if data != nil {
		if v, ok := data[name]; ok {
			return v, nil
		}
	}

	switch name {
	case "standard":
		return zincanalyzer.NewStandardAnalyzer(nil)
	case "simple":
		return analyzer.NewSimpleAnalyzer(), nil
	case "keyword":
		return analyzer.NewKeywordAnalyzer(), nil
	case "web":
		return analyzer.NewWebAnalyzer(), nil
	case "regexp", "pattern":
		return zincanalyzer.NewRegexpAnalyzer(nil)
	case "stop":
		return zincanalyzer.NewStopAnalyzer(nil)
	case "whitespace":
		return zincanalyzer.NewWhitespaceAnalyzer()
	case "gse_standard": // for Chinese support
		return chs.NewGseStandardAnalyzer(), nil
	case "gse_search": // for Chinese support
		return chs.NewGseSearchAnalyzer(), nil
		// language filters
	case "ar", "arabic":
		return ar.Analyzer(), nil
	case "cjk": // for Asia language
		return cjk.Analyzer(), nil
	case "ckb", "sorani":
		return ckb.Analyzer(), nil
	case "da", "danish":
		return da.Analyzer(), nil
	case "de", "german":
		return de.Analyzer(), nil
	case "en", "english":
		return en.NewAnalyzer(), nil
	case "es", "spanish":
		return es.Analyzer(), nil
	case "fa", "persian":
		return fa.Analyzer(), nil
	case "fi", "finnish":
		return fi.Analyzer(), nil
	case "fr", "french":
		return fr.Analyzer(), nil
	case "hi", "hindi":
		return hi.Analyzer(), nil
	case "hu", "hungarian":
		return hu.Analyzer(), nil
	case "it", "italian":
		return it.Analyzer(), nil
	case "nl", "dutch":
		return nl.Analyzer(), nil
	case "no", "norwegian":
		return no.Analyzer(), nil
	case "pt", "portuguese":
		return pt.Analyzer(), nil
	case "ro", "romanian":
		return ro.Analyzer(), nil
	case "ru", "russian":
		return ru.Analyzer(), nil
	case "sv", "swedish":
		return sv.Analyzer(), nil
	case "tr", "turkish":
		return tr.Analyzer(), nil
	default:
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] unknown analyzer [%s]", name))
	}
}

// QueryAnalyzerForField returns the analyzer and searchAnalyzer for the given field.
func QueryAnalyzerForField(data map[string]*analysis.Analyzer, mappings *meta.Mappings, field string) (*analysis.Analyzer, *analysis.Analyzer) {
	if field == "" {
		return nil, nil
	}

	analyzerName := ""
	searchAnalyzerName := ""
	if mappings != nil && mappings.Len() > 0 {
		if v, ok := mappings.GetProperty(field); ok {
			if v.Type != "text" {
				return nil, nil
			}
			if v.Analyzer != "" {
				analyzerName = v.Analyzer
			}
			if v.SearchAnalyzer != "" {
				searchAnalyzerName = v.SearchAnalyzer
			}
		}
	}

	analyzer, _ := QueryAnalyzer(data, analyzerName)
	searchAnalyzer, _ := QueryAnalyzer(data, searchAnalyzerName)

	return analyzer, searchAnalyzer
}
