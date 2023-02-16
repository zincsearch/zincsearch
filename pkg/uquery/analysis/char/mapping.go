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

package char

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge/analysis"

	zincchar "github.com/zinclabs/zincsearch/pkg/bluge/analysis/char"
	"github.com/zinclabs/zincsearch/pkg/errors"
	"github.com/zinclabs/zincsearch/pkg/zutils"
)

func NewMappingCharFilter(options interface{}) (analysis.CharFilter, error) {
	mappings, err := zutils.GetStringSliceFromMap(options, "mappings")
	if err != nil || len(mappings) == 0 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[char_filter] mapping option [mappings] should be exists")
	}
	for _, mapping := range mappings {
		if !strings.Contains(mapping, " => ") {
			return nil, errors.New(errors.ErrorTypeRuntimeException, fmt.Sprintf("[char_filter] mapping option [mappings] Invalid Mapping Rule: [%s], should be [old => new]", mapping))
		}
	}

	return zincchar.NewMappingCharFilter(mappings), nil
}
