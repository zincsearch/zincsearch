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

	zinctokenizer "github.com/zinclabs/zincsearch/pkg/bluge/analysis/tokenizer"
	"github.com/zinclabs/zincsearch/pkg/zutils"
)

func NewPathHierarchyTokenizer(options interface{}) (analysis.Tokenizer, error) {
	delimiter, _ := zutils.GetStringFromMap(options, "delimiter")
	if len(delimiter) == 0 {
		delimiter = "/"
	}
	replacement, _ := zutils.GetStringFromMap(options, "replacement")
	if len(replacement) == 0 {
		replacement = delimiter
	}
	skip, _ := zutils.GetFloatFromMap(options, "skip")
	return zinctokenizer.NewPathHierarchyTokenizer(delimiter[0], replacement[0], int(skip)), nil
}
