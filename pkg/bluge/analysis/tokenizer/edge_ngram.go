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
	"bytes"
	"unicode/utf8"

	"github.com/blugelabs/bluge/analysis"
)

type EdgeNgramTokenizer struct {
	minLength  int
	maxLength  int
	tokenChars []string
}

func NewEdgeNgramTokenizer(minLength, maxLength int, tokenChars []string) *EdgeNgramTokenizer {
	return &EdgeNgramTokenizer{
		minLength:  minLength,
		maxLength:  maxLength,
		tokenChars: tokenChars,
	}
}

func (t *EdgeNgramTokenizer) Tokenize(input []byte) analysis.TokenStream {
	n := utf8.RuneCount(input)
	runes := bytes.Runes(input)
	start := 0
	var byteStart = start
	rv := make(analysis.TokenStream, 0, n)
	for i := 1; i <= n; i++ {
		if i-start >= t.minLength {
			valid := true
			if len(t.tokenChars) > 0 {
				for _, c := range string(runes[start:i]) {
					if !t.isChar(c) {
						valid = false
						break
					}
				}
			}
			if valid {
				var term = analysis.BuildTermFromRunes(runes[start:i])
				rv = append(rv, &analysis.Token{
					Term:         term,
					PositionIncr: 1,
					Start:        byteStart,
					End:          byteStart + len(term),
					Type:         analysis.AlphaNumeric,
				})
			} else {
				byteStart = byteStart + utf8.RuneLen(runes[start])
				start = i
				continue
			}
		}

		if i-start == t.maxLength {
			break
		}
	}

	return rv
}

func (t *EdgeNgramTokenizer) isChar(r rune) bool {
	var ok bool
	for _, char := range t.tokenChars {
		if ok = isChar(char, r); ok {
			return true
		}
	}

	return false
}
