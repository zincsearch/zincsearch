package token

import (
	"github.com/blugelabs/bluge/analysis"

	"github.com/prabhatsharma/zinc/pkg/bluge/analysis/token"
)

func NewTrimTokenFilter() (analysis.TokenFilter, error) {
	return token.NewTrimTokenFilter(), nil
}
