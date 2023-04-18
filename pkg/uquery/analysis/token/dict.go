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

package token

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func NewDictTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	words, err := zutils.GetStringSliceFromMap(options, "words")
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, "[token_filter] dict option [words] should be an array of string")
	}
	dict := analysis.NewTokenMap()
	for _, word := range words {
		dict.AddToken(word)
	}

	minWordSize, _ := zutils.GetFloatFromMap(options, "min_word_size")
	minSubWordSize, _ := zutils.GetFloatFromMap(options, "min_sub_word_size")
	maxSubWordSize, _ := zutils.GetFloatFromMap(options, "max_sub_word_size")
	onlyLongestMatch, _ := zutils.GetBoolFromMap(options, "only_longest_match")
	if minWordSize == 0 {
		minWordSize = 5
	}
	if minSubWordSize == 0 {
		minSubWordSize = 2
	}
	if maxSubWordSize == 0 {
		maxSubWordSize = 15
	}
	return token.NewDictionaryCompoundFilter(dict, int(minWordSize), int(minSubWordSize), int(maxSubWordSize), onlyLongestMatch), nil
}
