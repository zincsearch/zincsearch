// Copyright 2022 Zinc Labs Inc. and Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"github.com/blugelabs/bluge"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/core"
)

func DeleteUser(userID string) bool {
	bdoc := bluge.NewDocument(userID)
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil))
	usersIndexWriter := core.ZINC_SYSTEM_INDEX_LIST["_users"].Writer
	err := usersIndexWriter.Delete(bdoc.ID())
	if err != nil {
		log.Printf("error deleting user: %s", err.Error())
		return false
	}

	// delete cache
	delete(ZINC_CACHED_USERS, userID)

	return true
}
