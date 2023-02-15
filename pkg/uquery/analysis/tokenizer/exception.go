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
	"regexp"
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/tokenizer"

	"github.com/zinclabs/zincsearch/pkg/errors"
	"github.com/zinclabs/zincsearch/pkg/zutils"
)

func NewExceptionTokenizer(options interface{}) (analysis.Tokenizer, error) {
	patterns, err := zutils.GetStringSliceFromMap(options, "patterns")
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, "[tokenizer] exception option [patterns] should be an array of string")
	}

	pattern := strings.Join(patterns, "|")
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[tokenizer] exception option [patterns] compile error: %s", err.Error()))
	}

	return tokenizer.NewExceptionsTokenizer(r, tokenizer.NewUnicodeTokenizer()), nil
}
