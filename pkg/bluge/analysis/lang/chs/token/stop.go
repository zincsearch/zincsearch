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
	"github.com/go-ego/gse"
)

type StopTokenFilter struct {
	seg *gse.Segmenter
}

func NewStopTokenFilter(seg *gse.Segmenter, stopwords []string) *StopTokenFilter {
	if len(stopwords) > 0 {
		for _, word := range stopwords {
			seg.AddStop(word)
		}
	}
	return &StopTokenFilter{seg}
}

func (f *StopTokenFilter) Filter(input analysis.TokenStream) analysis.TokenStream {
	var j, skipped int
	for _, token := range input {
		if !f.seg.IsStop(string(token.Term)) {
			token.PositionIncr += skipped
			skipped = 0
			input[j] = token
			j++
		} else {
			skipped += token.PositionIncr
		}
	}

	return input[:j]
}
