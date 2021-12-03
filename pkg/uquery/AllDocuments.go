package uquery

import (
	"github.com/blugelabs/bluge"
)

func AllDocuments() (bluge.SearchRequest, error) {

	query := bluge.NewMatchAllQuery()

	searchRequest := bluge.NewTopNSearch(1000, query)

	return searchRequest, nil

}
