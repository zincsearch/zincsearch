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

package analyzer

import (
	"fmt"
	"regexp"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"
	"github.com/blugelabs/bluge/analysis/tokenizer"

	zinctoken "github.com/zincsearch/zincsearch/pkg/bluge/analysis/token"
	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func NewRegexpAnalyzer(options interface{}) (*analysis.Analyzer, error) {
	pattern, _ := zutils.GetStringFromMap(options, "pattern")
	if pattern == "" {
		pattern = "\\w+"
	}
	lowerCase, err := zutils.GetBoolFromMap(options, "lowercase")
	if err != nil {
		lowerCase = true
	}
	stopwords, _ := zutils.GetStringSliceFromMap(options, "stopwords")

	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[analyzer] regexp option [pattern] compile error: %s", err.Error()))
	}

	ana := &analysis.Analyzer{Tokenizer: tokenizer.NewRegexpTokenizer(r)}
	if lowerCase {
		ana.TokenFilters = append(ana.TokenFilters, token.NewLowerCaseFilter())
	}

	if len(stopwords) > 0 {
		ana.TokenFilters = append(ana.TokenFilters, zinctoken.NewStopTokenFilter(stopwords))
	}

	return ana, nil
}
