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

	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func NewShingleTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	min, _ := zutils.GetFloatFromMap(options, "min_shingle_size")
	max, _ := zutils.GetFloatFromMap(options, "max_shingle_size")
	sep, _ := zutils.GetStringFromMap(options, "token_separator")
	fill, _ := zutils.GetStringFromMap(options, "filler_token")
	outputOriginalBool := true
	outputOriginal, _ := zutils.GetAnyFromMap(options, "output_original")
	if outputOriginal != nil {
		outputOriginalBool = outputOriginal.(bool)
	}
	if min == 0 {
		min = 2
	}
	if max == 0 {
		max = 2
	}
	if sep == "" {
		sep = " "
	}
	if fill == "" {
		fill = "_"
	}
	return token.NewShingleFilter(int(min), int(max), outputOriginalBool, sep, fill), nil
}
