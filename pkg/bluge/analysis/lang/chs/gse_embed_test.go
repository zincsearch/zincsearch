//go:build !ne

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

package chs

import (
	"testing"

	"github.com/go-ego/gse"
	"github.com/stretchr/testify/assert"
)

// "科幻片" is only in gse's embedded big dict (small dict splits it), so BIG must really load it here.
func TestLoadDict_BigEmbedded(t *testing.T) {
	old := seg
	seg = new(gse.Segmenter)
	t.Cleanup(func() { seg = old })

	assert.True(t, bigDictEmbedded)
	loadDict(true, false, "BIG")
	got := NewGseStandardTokenizer().Tokenize([]byte("科幻片"))
	assert.Equal(t, "[科幻片]", collectToken(got))
}
