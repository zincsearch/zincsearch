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
	"github.com/blugelabs/bluge/analysis/lang/ga"
	"github.com/blugelabs/bluge/analysis/lang/hi"
	"github.com/blugelabs/bluge/analysis/lang/hu"
	"github.com/blugelabs/bluge/analysis/lang/in"
	"github.com/blugelabs/bluge/analysis/lang/it"
	"github.com/blugelabs/bluge/analysis/lang/nl"
	"github.com/blugelabs/bluge/analysis/lang/no"
	"github.com/blugelabs/bluge/analysis/lang/pt"
	"github.com/blugelabs/bluge/analysis/lang/ro"
	"github.com/blugelabs/bluge/analysis/lang/ru"
	"github.com/blugelabs/bluge/analysis/lang/sv"
	"github.com/blugelabs/bluge/analysis/lang/tr"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/zinclabs/zincsearch/pkg/bluge/analysis/lang/chs"
	"github.com/zinclabs/zincsearch/pkg/errors"
	zinctoken "github.com/zinclabs/zincsearch/pkg/uquery/analysis/token"
	"github.com/zinclabs/zincsearch/pkg/zutils"
)

func RequestTokenFilter(data map[string]interface{}) (map[string]analysis.TokenFilter, error) {
	if data == nil {
		return nil, nil
	}

	filters := make(map[string]analysis.TokenFilter)
	for name, options := range data {
		typ, err := zutils.GetStringFromMap(options, "type")
		if err != nil {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[token_filter] %s option [%s] should be exists", name, "type"))
		}
		filter, err := RequestTokenFilterSingle(typ, options)
		if err != nil {
			return nil, err
		}
		filters[name] = filter
	}

	return filters, nil
}

func RequestTokenFilterSlice(data []interface{}) ([]analysis.TokenFilter, error) {
	if data == nil {
		return nil, nil
	}

	filters := make([]analysis.TokenFilter, 0, len(data))
	for _, options := range data {
		var err error
		var filter analysis.TokenFilter
		switch v := options.(type) {
		case string:
			filter, err = RequestTokenFilterSingle(v, nil)
		case map[string]interface{}:
			var typ string
			typ, err = zutils.GetStringFromMap(options, "type")
			if err != nil {
				return nil, errors.New(errors.ErrorTypeParsingException, "[token_filter] option [type] should be exists")
			}
			filter, err = RequestTokenFilterSingle(typ, options)
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, "[token_filter] option should be string or object")
		}
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}

	return filters, nil
}

func RequestTokenFilterSingle(name string, options interface{}) (analysis.TokenFilter, error) {
	name = strings.ToLower(name)
	switch name {
	case "apostrophe":
		return token.NewApostropheFilter(), nil
	case "camel_case", "camelcase":
		return token.NewCamelCaseFilter(), nil
	case "dict":
		return zinctoken.NewDictTokenFilter(options)
	case "edge_ngram":
		return zinctoken.NewEdgeNgramTokenFilter(options)
	case "elision":
		return zinctoken.NewElisionTokenFilter(options)
	case "keyword", "keyword_marker":
		return zinctoken.NewKeywordTokenFilter(options)
	case "length":
		return zinctoken.NewLengthTokenFilter(options)
	case "lower_case", "lowercase":
		return token.NewLowerCaseFilter(), nil
	case "ngram":
		return zinctoken.NewNgramTokenFilter(options)
	case "porter", "stemmer":
		return token.NewPorterStemmer(), nil
	case "reverse":
		return token.NewReverseFilter(), nil
	case "regexp", "pattern_replace":
		return zinctoken.NewRegexpTokenFilter(options)
	case "shingle":
		return zinctoken.NewShingleTokenFilter(options)
	case "trim":
		return zinctoken.NewTrimTokenFilter()
	case "stop":
		return zinctoken.NewStopTokenFilter(options)
	case "truncate":
		return zinctoken.NewTruncateTokenFilter(options)
	case "unicodenorm":
		return zinctoken.NewUnicodenormTokenFilter(options)
	case "unique":
		return token.NewUniqueTermFilter(), nil
	case "upper_case", "uppercase":
		return zinctoken.NewUpperCaseTokenFilter()
	case "gse_stop":
		return chs.NewGseStopTokenFilter(), nil
		// language filters
	case "ar_normalization", "arabic_normalization":
		return ar.NormalizeFilter(), nil
	case "ar_stemmer", "arabic_stemmer":
		return ar.StemmerFilter(), nil
	case "cjk_bigram":
		return cjk.NewBigramFilter(false), nil
	case "cjk_width":
		return cjk.NewWidthFilter(), nil
	case "ckb_normalization", "sorani_normalization":
		return ckb.NormalizeFilter(), nil
	case "ckb_stemmer", "sorani_stemmer":
		return ckb.StemmerFilter(), nil
	case "da_stemmer", "danish_stemmer":
		return da.StemmerFilter(), nil
	case "de_normalization", "german_normalization":
		return de.NormalizeFilter(), nil
	case "de_stemmer", "german_stemmer":
		return de.StemmerFilter(), nil
	case "de_light_stemmer", "german_light_stemmer":
		return de.LightStemmerFilter(), nil
	case "en_possessive_stemmer", "english_possessive_stemmer":
		return en.NewPossessiveFilter(), nil
	case "en_stemmer", "english_stemmer":
		return en.StemmerFilter(), nil
	case "es_stemmer", "spanish_stemmer":
		return es.StemmerFilter(), nil
	case "es_light_stemmer", "spanish_light_stemmer":
		return es.LightStemmerFilter(), nil
	case "fa_normalization", "persian_normalization":
		return fa.NormalizeFilter(), nil
	case "fi_stemmer", "finnish_stemmer":
		return fi.StemmerFilter(), nil
	case "fr_elision", "french_elision":
		return fr.ElisionFilter(), nil
	case "fr_stemmer", "french_stemmer":
		return fr.StemmerFilter(), nil
	case "fr_light_stemmer", "french_light_stemmer":
		return fr.LightStemmerFilter(), nil
	case "fr_minimal_stemmer", "french_minimal_stemmer":
		return fr.MinimalStemmerFilter(), nil
	case "ga_elision", "irish_elision":
		return ga.ElisionFilter(), nil
	case "hi_normalization", "hindi_normalization":
		return hi.NormalizeFilter(), nil
	case "hi_stemmer", "hindi_stemmer":
		return hi.StemmerFilter(), nil
	case "hu_stemmer", "hungarian_stemmer":
		return hu.StemmerFilter(), nil
	case "in_normalization", "indic_normalization":
		return in.NormalizeFilter(), nil
	case "it_elision", "italian_elision":
		return it.ElisionFilter(), nil
	case "it_stemmer", "italian_stemmer":
		return it.StemmerFilter(), nil
	case "it_light_stemmer", "italian_light_stemmer":
		return it.LightStemmerFilter(), nil
	case "nl_stemmer", "dutch_stemmer":
		return nl.StemmerFilter(), nil
	case "no_stemmer", "norwegian_stemmer":
		return no.StemmerFilter(), nil
	case "pt_light_stemmer", "portuguese_stemmer", "portuguese_light_stemmer":
		return pt.LightStemmerFilter(), nil
	case "ro_stemmer", "romanian_stemmer":
		return ro.StemmerFilter(), nil
	case "ru_stemmer", "russian_stemmer":
		return ru.StemmerFilter(), nil
	case "sv_stemmer", "swedish_stemmer":
		return sv.StemmerFilter(), nil
	case "tr_stemmer", "turkish_stemmer":
		return tr.StemmerFilter(), nil
	default:
		return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[token_filter] unknown token filter [%s]", name))
	}
}
