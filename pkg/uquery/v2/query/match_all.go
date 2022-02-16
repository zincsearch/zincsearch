package query

import "github.com/blugelabs/bluge"

func MatchAllQuery(query map[string]interface{}) (bluge.Query, error) {
	return bluge.NewMatchAllQuery(), nil
}
