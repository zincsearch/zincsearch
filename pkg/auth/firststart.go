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
	"os"

	"github.com/rs/zerolog/log"
)

var ZINC_CACHED_USERS map[string]SimpleUser

func init() {
	// initialize cache
	ZINC_CACHED_USERS = make(map[string]SimpleUser)

	firstStart, err := IsFirstStart()
	if err != nil {
		log.Print(err)
	}
	if firstStart {
		// create default user from environment variable
		adminUser := os.Getenv("ZINC_FIRST_ADMIN_USER")
		adminPassword := os.Getenv("ZINC_FIRST_ADMIN_PASSWORD")

		if adminUser == "" || adminPassword == "" {
			log.Fatal().Msg("ZINC_FIRST_ADMIN_USER and ZINC_FIRST_ADMIN_PASSWORD must be set on first start. You should also change the credentials after first login.")
		}
		CreateUser(adminUser, adminUser, adminPassword, "admin")
	}
}

func IsFirstStart() (bool, error) {
	userList, err := GetAllUsersWorker()
	if err != nil {
		return true, err
	}

	if userList.Hits.Total.Value == 0 {
		return true, nil
	}

	// cache users
	for _, user := range userList.Hits.Hits {
		ZINC_CACHED_USERS[user.ID] = user.Source.(SimpleUser)
	}

	return false, nil
}
