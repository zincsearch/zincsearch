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
	"errors"
	"os"

	"github.com/rs/zerolog/log"
)

func init() {
	// init first start
	firstStart, err := isFirstStart()
	if err != nil {
		log.Print(err)
	}
	if firstStart {
		if err := initFirstUser(); err != nil {
			log.Fatal().Err(err).Msg("init first user")
		}
	}
	if err := initPermissionCache(); err != nil {
		log.Print(err)
	}
}

func isFirstStart() (bool, error) {
	users, err := GetUsers()
	if err != nil {
		return true, err
	}

	if len(users) == 0 {
		return true, nil
	}

	for _, user := range users {
		ZINC_CACHED_USERS.Set(user.ID, user)
	}

	return false, nil
}

func initPermissionCache() error {
	roles, err := GetRoles()
	if err != nil {
		return err
	}

	for _, role := range roles {
		ZINC_CACHED_PERMISSIONS.Set(role.ID, strArrayToMap(role.Permission))
	}

	return nil
}

func initFirstUser() error {
	// create default user from environment variable
	adminUser := os.Getenv("ZINC_FIRST_ADMIN_USER")
	adminPassword := os.Getenv("ZINC_FIRST_ADMIN_PASSWORD")
	if adminUser == "" || adminPassword == "" {
		return errors.New("ZINC_FIRST_ADMIN_USER and ZINC_FIRST_ADMIN_PASSWORD must be set on first start. You should also change the credentials after first login")
	}

	_, err := CreateUser(adminUser, adminUser, adminPassword, "admin")

	return err
}
