package query

import "github.com/blugelabs/bluge"

func MatchNoneQuery() (bluge.Query, error) {
	return bluge.NewMatchNoneQuery(), nil
}
