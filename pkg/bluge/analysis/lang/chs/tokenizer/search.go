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
	"github.com/go-ego/gse"

	"github.com/zincsearch/zincsearch/pkg/config"
)

type SearchTokenizer struct {
	seg *gse.Segmenter
}

func NewSearchTokenizer(seg *gse.Segmenter) *SearchTokenizer {
	return &SearchTokenizer{seg}
}

func (t *SearchTokenizer) Tokenize(input []byte) analysis.TokenStream {
	result := make(analysis.TokenStream, 0, len(input))
	text := string(input)
	search := t.seg.CutSearch(text, config.Global.Plugin.GSE.EnableHMM)
	tokens := t.seg.Analyze(search, text)
	var start, positionIncr int
	for _, token := range tokens {
		positionIncr = 1
		if start == token.Start {
			positionIncr = 0
		}
		start = token.Start

		typ := analysis.Ideographic
		alphaNumeric := true
		for _, r := range token.Text {
			if r < 32 || r > 126 {
				alphaNumeric = false
				break
			}
		}
		if alphaNumeric {
			typ = analysis.AlphaNumeric
		}

		result = append(result, &analysis.Token{
			Term:         []byte(token.Text),
			Start:        token.Start,
			End:          token.End,
			PositionIncr: positionIncr,
			Type:         typ,
		})
	}
	return result
}
