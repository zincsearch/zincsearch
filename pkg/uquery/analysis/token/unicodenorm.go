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
	"fmt"
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"
	"golang.org/x/text/unicode/norm"

	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func NewUnicodenormTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	form, _ := zutils.GetStringFromMap(options, "form")
	form = strings.ToUpper(form)
	switch form {
	case "NFC":
		return token.NewUnicodeNormalizeFilter(norm.NFC), nil
	case "NFD":
		return token.NewUnicodeNormalizeFilter(norm.NFD), nil
	case "NFKC":
		return token.NewUnicodeNormalizeFilter(norm.NFKC), nil
	case "NFKD":
		return token.NewUnicodeNormalizeFilter(norm.NFKD), nil
	default:
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[token_filter] unicodenorm doesn't support form [%s]", form))
	}
}
