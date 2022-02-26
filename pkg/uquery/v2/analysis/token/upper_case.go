package token

import (
	"github.com/blugelabs/bluge/analysis"

	"github.com/prabhatsharma/zinc/pkg/bluge/analysis/token"
)

func NewUpperCaseTokenFilter() (analysis.TokenFilter, error) {
	return token.NewUpperCaseTokenFilter(), nil
}
