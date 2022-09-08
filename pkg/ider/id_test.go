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

package ider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	got := Generate()
	assert.NotEmpty(t, got)
}

func TestNewNode(t *testing.T) {
	for i := 1023; i < 1026; i++ {
		node, err := NewNode(i)
		assert.NoError(t, err)
		assert.NotNil(t, node)
		id := node.Generate()
		assert.NotEmpty(t, id)
	}
}
