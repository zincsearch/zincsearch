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

package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAllUsersWorker(t *testing.T) {
	t.Run("prepare", func(t *testing.T) {
		u, err := CreateUser("test", "test", "test", "admin")
		assert.NoError(t, err)
		assert.NotNil(t, u)
	})

	t.Run("get all users", func(t *testing.T) {
		// wait for _users prepared
		time.Sleep(time.Second)
		got, err := GetUsers()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(got), 1)
	})
}
