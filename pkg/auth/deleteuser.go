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
