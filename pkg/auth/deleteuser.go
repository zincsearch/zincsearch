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
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/metadata"
)

func DeleteUser(userID string) bool {
	err := metadata.User.Delete(userID)
	if err != nil {
		log.Error().Err(err).Msg("error deleting user")
		return false
	}

	// delete cache
	delete(ZINC_CACHED_USERS, userID)

	return true
}
