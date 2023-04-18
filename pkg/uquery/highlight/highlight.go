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

package highlight

import (
	"github.com/zincsearch/zincsearch/pkg/meta"
)

func Request(highlight *meta.Highlight) error {
	if len(highlight.Fields) == 0 {
		return nil
	}

	if highlight.NumberOfFragments == 0 {
		highlight.NumberOfFragments = 3
	}
	for _, field := range highlight.Fields {
		if field.FragmentSize == 0 && highlight.FragmentSize > 0 {
			field.FragmentSize = highlight.FragmentSize
		}
		if field.NumberOfFragments == 0 && highlight.NumberOfFragments > 0 {
			field.NumberOfFragments = highlight.NumberOfFragments
		}
	}

	return nil
}
