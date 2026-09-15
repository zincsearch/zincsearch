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
	"time"

	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/metadata"
)

var ErrInvalidCredentials = errors.New(errors.ErrorTypeInvalidArgument, "invalid credentials")

var accountUserStore interface {
	Set(string, meta.User) error
} = metadata.User

// UpdateAccount lets a user change their own name and/or password after
// proving the current password. Empty name or newPassword keeps the current value.
func UpdateAccount(id, password, name, newPassword string) (*meta.User, error) {
	if name == "" && newPassword == "" {
		return nil, errors.New(errors.ErrorTypeInvalidArgument, "name or new password is required")
	}
	user, ok := VerifyCredentials(id, password)
	if !ok {
		return nil, ErrInvalidCredentials
	}

	updated := *user
	if name != "" {
		updated.Name = name
	}
	if newPassword != "" {
		updated.Salt = GenerateSalt()
		updated.Password = GeneratePassword(newPassword, updated.Salt)
	}
	updated.UpdatedAt = time.Now()
	if err := accountUserStore.Set(updated.ID, updated); err != nil {
		return nil, err
	}
	ZINC_CACHED_USERS.Set(updated.ID, &updated)
	return &updated, nil
}
