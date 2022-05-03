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
	"unicode"

	"github.com/blugelabs/bluge/analysis/tokenizer"
)

type IsTokenRune func(r rune) bool

type CharGroupTokenizer struct {
	chars []string
}

func NewCharGroupTokenizer(chars []string) *tokenizer.CharacterTokenizer {
	t := new(CharGroupTokenizer)
	for _, char := range chars {
		if len(char) == 0 {
			continue
		}
		t.chars = append(t.chars, char)
	}
	return tokenizer.NewCharacterTokenizer(t.isChar)
}

func (t *CharGroupTokenizer) isChar(r rune) bool {
	var ok bool
	for _, char := range t.chars {
		if ok = isChar(char, r); ok {
			return false
		}
	}

	return true
}

func isChar(char string, r rune) bool {
	ok := false
	switch char {
	case "graphic":
		ok = unicode.IsGraphic(r)
	case "print":
		ok = unicode.IsPrint(r)
	case "control":
		ok = unicode.IsControl(r)
	case "letter":
		ok = unicode.IsLetter(r)
	case "mark":
		ok = unicode.IsMark(r)
	case "number", "digit":
		ok = unicode.IsNumber(r)
	case "punct", "punctuation":
		ok = unicode.IsPunct(r)
	case "space", "whitespace", "white_space":
		ok = unicode.IsSpace(r)
	case "symbol":
		ok = unicode.IsSymbol(r)
	default:
		for _, c := range char {
			if r == c {
				ok = true
				break
			}
		}
	}

	return ok
}
