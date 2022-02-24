package tokenizer

import (
	"github.com/blugelabs/bluge/analysis"
)

type EdgeNgramTokenizer struct {
	minLength int
	maxLength int
}

func NewEdgeNgramTokenizer(minLength, maxLength int) *EdgeNgramTokenizer {
	return &EdgeNgramTokenizer{
		minLength: minLength,
		maxLength: maxLength,
	}
}

func (t *EdgeNgramTokenizer) Tokenize(input []byte) analysis.TokenStream {
	n := len(input)
	if n > t.maxLength {
		n = t.maxLength
	}

	rv := make(analysis.TokenStream, 0, n)
	for i := t.minLength; i <= n; i++ {
		rv = append(rv, &analysis.Token{
			Term:         input[:i],
			PositionIncr: 1,
			Start:        0,
			End:          i,
			Type:         analysis.AlphaNumeric,
		})
	}

	return rv
}
