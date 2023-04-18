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

func NewElisionTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	articles, err := zutils.GetStringSliceFromMap(options, "articles")
	if err != nil {
		articles = []string{"l", "m", "t", "qu", "n", "s", "j", "d", "c", "jusqu", "quoiqu", "lorsqu", "puisqu"}
	}
	dict := analysis.NewTokenMap()
	for _, word := range articles {
		dict.AddToken(word)
	}
	return token.NewElisionFilter(dict), nil
}
