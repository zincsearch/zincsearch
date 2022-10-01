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
	"strings"

	"github.com/zinclabs/zinc/pkg/meta"
)

func VerifyCredentials(userID, password string) (*meta.User, bool) {
	userID = strings.ToLower(userID)
	user, ok := ZINC_CACHED_USERS.Get(userID)
	if !ok {
		return user, false
	}

	incomingEncryptedPassword := GeneratePassword(password, user.Salt)
	if incomingEncryptedPassword == user.Password {
		return user, true
	}

	return nil, false
}

func VerifyRoleHasPermission(roleId, permission string) bool {
	roleId = strings.ToLower(roleId)
	if roleId == "admin" {
		return true
	}
	pm, ok := ZINC_CACHED_PERMISSIONS.Get(roleId)
	if !ok {
		return false
	}
	_, ok = pm[permission]
	return ok
}
