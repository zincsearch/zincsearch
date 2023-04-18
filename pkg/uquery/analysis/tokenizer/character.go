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
	"fmt"
	"unicode"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/tokenizer"

	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func NewCharacterTokenizer(options interface{}) (analysis.Tokenizer, error) {
	char, _ := zutils.GetStringFromMap(options, "char")
	switch char {
	case "graphic":
		return tokenizer.NewCharacterTokenizer(unicode.IsGraphic), nil
	case "print":
		return tokenizer.NewCharacterTokenizer(unicode.IsPrint), nil
	case "control":
		return tokenizer.NewCharacterTokenizer(unicode.IsControl), nil
	case "letter":
		return tokenizer.NewCharacterTokenizer(unicode.IsLetter), nil
	case "mark":
		return tokenizer.NewCharacterTokenizer(unicode.IsMark), nil
	case "number", "digit":
		return tokenizer.NewCharacterTokenizer(unicode.IsNumber), nil
	case "punct", "punctuation":
		return tokenizer.NewCharacterTokenizer(unicode.IsPunct), nil
	case "space", "whitespace", "white_space":
		return tokenizer.NewCharacterTokenizer(unicode.IsSpace), nil
	case "symbol":
		return tokenizer.NewCharacterTokenizer(unicode.IsSymbol), nil
	default:
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[tokenizer] character doesn't support char [%s]", char))
	}
}
