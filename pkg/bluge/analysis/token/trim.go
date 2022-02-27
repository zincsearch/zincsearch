package token

import (
	"bytes"

	"github.com/blugelabs/bluge/analysis"
)

type TrimTokenFilter struct{}

func NewTrimTokenFilter() *TrimTokenFilter {
	return &TrimTokenFilter{}
}

func (t *TrimTokenFilter) Filter(input analysis.TokenStream) analysis.TokenStream {
	for _, token := range input {
		token.Term = bytes.TrimSpace(token.Term)
	}

	return input
}
