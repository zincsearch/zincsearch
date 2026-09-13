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

package zincsearch

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFrontendAssets(t *testing.T) {
	t.Run("embed::GetFrontendAssets", func(t *testing.T) {
		f, err := GetFrontendAssets()
		require.NoError(t, err)
		require.NotNil(t, f)
		t.Run("index.html", func(t *testing.T) {
			content, err := fs.ReadFile(f, "index.html")
			require.NoError(t, err)
			assert.NotEmpty(t, content)

			expected, err := embedFrontend.ReadFile("frontend/dist/index.html")
			require.NoError(t, err)
			assert.Equal(t, expected, content)
		})
	})
}
