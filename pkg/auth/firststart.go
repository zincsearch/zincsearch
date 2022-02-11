package auth

import (
	"log"
	"os"
)

func init() {
	firstStart, err := IsFirstStart()
	if err != nil {
		log.Println(err)
	}
	if firstStart {
		// create default user from environment variable
		adminUser := os.Getenv("ZINC_FIRST_ADMIN_USER")
		adminPassword := os.Getenv("ZINC_FIRST_ADMIN_PASSWORD")

		if adminUser == "" || adminPassword == "" {
			log.Fatal("FIRST_ADMIN_USER and FIRST_ADMIN_PASSWORD must be set on first start. You should also change the credentials after first login.")
		}
		CreateUser(adminUser, "", adminPassword, "admin")
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

	return false, nil

}
