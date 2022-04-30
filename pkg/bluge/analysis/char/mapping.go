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
	"bytes"
	"strings"
)

type MappingCharFilter struct {
	old [][]byte
	new [][]byte
}

func NewMappingCharFilter(mappings []string) *MappingCharFilter {
	m := &MappingCharFilter{}
	for _, field := range mappings {
		field := strings.Split(field, " => ")
		if len(field) != 2 {
			continue
		}
		m.old = append(m.old, []byte(field[0]))
		m.new = append(m.new, []byte(field[1]))
	}

	return m
}

func (t *MappingCharFilter) Filter(input []byte) []byte {
	for i := 0; i < len(t.old); i++ {
		input = []byte(bytes.ReplaceAll(input, t.old[i], t.new[i]))
	}
	return input
}
