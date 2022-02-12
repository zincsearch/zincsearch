package parser

import "github.com/blugelabs/bluge"

func SimpleQueryStringQuery(query map[string]interface{}) (bluge.Query, error) {
	return QueryStringQuery(query)
}
