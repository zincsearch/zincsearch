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

	"github.com/zincsearch/zincsearch/pkg/zutils/base62"
)

func TestGenerate(t *testing.T) {
	got := Generate()
	assert.NotEmpty(t, got)
}

func TestNewNode(t *testing.T) {
	for _, i := range []int{0, 1, maxNodeID} {
		node, err := NewNode(i)
		assert.NoError(t, err)
		assert.NotNil(t, node)
		id := node.Generate()
		assert.NotEmpty(t, id)
	}
	for _, i := range []int{-1, maxNodeID + 1, 2048} {
		node, err := NewNode(i)
		assert.Error(t, err, "node id %d must be rejected", i)
		assert.Nil(t, node)
	}
}

func TestGenerate_Unique(t *testing.T) {
	node, err := NewNode(1)
	assert.NoError(t, err)
	const n = 10000
	seen := make(map[string]struct{}, n)
	prev := int64(0)
	for range n {
		raw, err := node.node.NextID()
		assert.NoError(t, err)
		assert.Greater(t, raw, prev, "ids must be strictly increasing")
		assert.Equal(t, int64(1), raw&1023, "low 10 bits carry the node id")
		prev = raw
		id := base62.Encode(raw)
		_, dup := seen[id]
		assert.False(t, dup, "duplicate id %s", id)
		seen[id] = struct{}{}
	}
}
