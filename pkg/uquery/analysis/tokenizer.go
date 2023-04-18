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
	"github.com/blugelabs/bluge/analysis/tokenizer"

	"github.com/zincsearch/zincsearch/pkg/bluge/analysis/lang/chs"
	"github.com/zincsearch/zincsearch/pkg/errors"
	zinctokenizer "github.com/zincsearch/zincsearch/pkg/uquery/analysis/tokenizer"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func RequestTokenizer(data map[string]interface{}) (map[string]analysis.Tokenizer, error) {
	if data == nil {
		return nil, nil
	}

	tokenizers := make(map[string]analysis.Tokenizer)
	for name, options := range data {
		typ, err := zutils.GetStringFromMap(options, "type")
		if err != nil {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[tokenizer] %s option [%s] should be exists", name, "type"))
		}
		zer, err := RequestTokenizerSingle(typ, options)
		if err != nil {
			return nil, err
		}
		tokenizers[name] = zer
	}

	return tokenizers, nil
}

func RequestTokenizerSlice(data []interface{}) ([]analysis.Tokenizer, error) {
	if data == nil {
		return nil, nil
	}

	tokenizers := make([]analysis.Tokenizer, 0, len(data))
	for _, typ := range data {
		typ, ok := typ.(string)
		if !ok {
			return nil, errors.New(errors.ErrorTypeParsingException, "[tokenizer] option should be string")
		}
		zer, err := RequestTokenizerSingle(typ, nil)
		if err != nil {
			return nil, err
		}
		tokenizers = append(tokenizers, zer)
	}

	return tokenizers, nil
}

func RequestTokenizerSingle(name string, options interface{}) (analysis.Tokenizer, error) {
	name = strings.ToLower(name)
	switch name {
	case "character":
		return zinctokenizer.NewCharacterTokenizer(options)
	case "char_group":
		return zinctokenizer.NewCharGroupTokenizer(options)
	case "edge_ngram":
		return zinctokenizer.NewEdgeNgramTokenizer(options)
	case "exception":
		return zinctokenizer.NewExceptionTokenizer(options)
	case "letter", "simple":
		return tokenizer.NewLetterTokenizer(), nil
	case "lower_case", "lowercase":
		return zinctokenizer.NewLowerCaseTokenizer()
	case "ngram":
		return zinctokenizer.NewNgramTokenizer(options)
	case "path_hierarchy":
		return zinctokenizer.NewPathHierarchyTokenizer(options)
	case "regexp", "pattern":
		return zinctokenizer.NewRegexpTokenizer(options)
	case "single", "keyword":
		return tokenizer.NewSingleTokenTokenizer(), nil
	case "unicode", "standard":
		return tokenizer.NewUnicodeTokenizer(), nil
	case "web":
		return tokenizer.NewWebTokenizer(), nil
	case "whitespace":
		return tokenizer.NewWhitespaceTokenizer(), nil
	case "gse_standard":
		return chs.NewGseStandardTokenizer(), nil
	case "gse_search":
		return chs.NewGseSearchTokenizer(), nil
	default:
		return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[tokenizer] unknown tokenizer [%s]", name))
	}
}
