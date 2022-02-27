package token

import (
	"regexp"

	"github.com/blugelabs/bluge/analysis"
)

type RegexpTokenFilter struct {
	r           *regexp.Regexp
	replacement []byte
}

func NewRegexpTokenFilter(r *regexp.Regexp, replacement []byte) *RegexpTokenFilter {
	return &RegexpTokenFilter{
		r:           r,
		replacement: replacement,
	}
}

func (t *RegexpTokenFilter) Filter(input analysis.TokenStream) analysis.TokenStream {
	for _, token := range input {
		token.Term = t.r.ReplaceAll(token.Term, t.replacement)
	}

	return input
}
