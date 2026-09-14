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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/pkg/meta"
)

func TestRequest(t *testing.T) {
	t.Run("no fields is a noop", func(t *testing.T) {
		h := &meta.Highlight{}
		assert.NoError(t, Request(h))
		assert.Equal(t, 0, h.NumberOfFragments)
	})

	t.Run("defaults propagate to fields", func(t *testing.T) {
		h := &meta.Highlight{
			FragmentSize: 100,
			Fields: map[string]*meta.Highlight{
				"title": {},
				"body":  {FragmentSize: 20, NumberOfFragments: 1},
			},
		}
		assert.NoError(t, Request(h))
		assert.Equal(t, 3, h.NumberOfFragments)
		assert.Equal(t, 100, h.Fields["title"].FragmentSize)
		assert.Equal(t, 3, h.Fields["title"].NumberOfFragments)
		// explicit per-field values are kept
		assert.Equal(t, 20, h.Fields["body"].FragmentSize)
		assert.Equal(t, 1, h.Fields["body"].NumberOfFragments)
	})

	t.Run("zero global fragment size leaves field untouched", func(t *testing.T) {
		h := &meta.Highlight{NumberOfFragments: 5, Fields: map[string]*meta.Highlight{"f": {}}}
		assert.NoError(t, Request(h))
		assert.Equal(t, 0, h.Fields["f"].FragmentSize)
		assert.Equal(t, 5, h.Fields["f"].NumberOfFragments)
	})
}
