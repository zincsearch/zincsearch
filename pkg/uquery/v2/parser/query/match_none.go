package query

import "github.com/blugelabs/bluge"

func MatchNoneQuery(query map[string]interface{}) (bluge.Query, error) {
	return bluge.NewMatchNoneQuery(), nil
}
