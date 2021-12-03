package core

import (
	"github.com/blugelabs/bluge"
)

// Nothing to handle in the error. If you can't load indexes then everything is broken.
var ZINC_INDEX_LIST, _ = LoadZincIndexes()

var ZINC_SYSTEM_INDEX_LIST, _ = LoadZincSystemIndexes()

type Index struct {
	Name          string            `json:"name"`
	Writer        *bluge.Writer     `json:"blugeWriter"`
	CachedMapping map[string]string `json:"mapping"`
	IndexType     string            `json:"index_type"` // "system" or "user"
}
