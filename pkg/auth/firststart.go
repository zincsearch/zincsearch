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
