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

func NewLengthTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	min, _ := zutils.GetFloatFromMap(options, "min")
	max, _ := zutils.GetFloatFromMap(options, "max")
	if min == 0 {
		min = 1
	}
	if max == 0 {
		max = 2
	}
	return token.NewLengthFilter(int(min), int(max)), nil
}
