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

package tokenizer

import (
	"github.com/blugelabs/bluge/analysis"

	zinctokenizer "github.com/zincsearch/zincsearch/pkg/bluge/analysis/tokenizer"
	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func NewEdgeNgramTokenizer(options interface{}) (analysis.Tokenizer, error) {
	min, _ := zutils.GetFloatFromMap(options, "min_gram")
	max, _ := zutils.GetFloatFromMap(options, "max_gram")
	tokenChars, _ := zutils.GetStringSliceFromMap(options, "token_chars")
	if min == 0 {
		min = 1
	}
	if max == 0 {
		max = 2
	}
	if min > max {
		return nil, errors.New(errors.ErrorTypeParsingException, "[tokenizer] edge_ngram option [min_gram] should be not greater than [max_gram]")
	}
	return zinctokenizer.NewEdgeNgramTokenizer(int(min), int(max), tokenChars), nil
}
