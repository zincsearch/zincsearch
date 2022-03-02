package query

import "github.com/blugelabs/bluge"

func MatchAllQuery() (bluge.Query, error) {
	return bluge.NewMatchAllQuery(), nil
}
