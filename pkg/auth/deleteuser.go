package auth

import (
	"github.com/blugelabs/bluge"
	"github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc/pkg/core"
)

func DeleteUser(userId string) bool {
	bdoc := bluge.NewDocument(userId)
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil))
	usersIndexWriter := core.ZINC_SYSTEM_INDEX_LIST["_users"].Writer
	err := usersIndexWriter.Delete(bdoc.ID())
	if err != nil {
		log.Printf("error deleting user: %v", err)
		return false
	}

	return true
}
