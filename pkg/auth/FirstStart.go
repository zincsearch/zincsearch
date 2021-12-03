package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func init() {
	firstStart, err := IsFirstStart()
	if err != nil {
		fmt.Println(err)
	}
	if firstStart {
		// create default user from environment variable
		adminUser := os.Getenv("FIRST_ADMIN_USER")
		adminPassword := os.Getenv("FIRST_ADMIN_PASSWORD")

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

	// Logger(userList)

	if userList.Hits.Total.Value == 0 {
		return true, nil
	}

	return false, nil

}

func Logger(m interface{}) {

	k1, _ := json.Marshal(m)

	var k2 map[string]interface{}
	json.Unmarshal(k1, &k2)
	k2["time"] = time.Now()

	k3, _ := json.Marshal(k2)
	fmt.Println(string(k3))
}
