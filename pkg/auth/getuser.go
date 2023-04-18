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
	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/metadata"
)

func GetUser(id string) (*meta.User, bool, error) {
	if id == "" {
		return nil, false, errors.New(errors.ErrorTypeInvalidArgument, "user id is required")
	}
	user, err := metadata.User.Get(id)
	if err != nil {
		return nil, false, err
	}
	return user, true, nil
}

func SetUser(id string, user meta.User) error {
	return metadata.User.Set(id, user)
}
