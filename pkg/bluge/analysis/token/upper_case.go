package token

import (
	"bytes"

	"github.com/blugelabs/bluge/analysis"
)

type UpperCaseTokenFilter struct{}

func NewUpperCaseTokenFilter() *UpperCaseTokenFilter {
	return &UpperCaseTokenFilter{}
}

func (t *UpperCaseTokenFilter) Filter(input analysis.TokenStream) analysis.TokenStream {
	for _, token := range input {
		token.Term = bytes.ToUpper(token.Term)
	}

	return input
}
