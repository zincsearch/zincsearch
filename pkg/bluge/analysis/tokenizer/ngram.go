package tokenizer

import (
	"github.com/blugelabs/bluge/analysis"
)

type NgramTokenizer struct {
	minLength int
	maxLength int
}

func NewNgramTokenizer(minLength, maxLength int) *NgramTokenizer {
	return &NgramTokenizer{
		minLength: minLength,
		maxLength: maxLength,
	}
}

func (t *NgramTokenizer) Tokenize(input []byte) analysis.TokenStream {
	n := len(input)
	start := 0
	rv := make(analysis.TokenStream, 0, n)
	for i := 1; i <= n; i++ {
		rv = append(rv, &analysis.Token{
			Term:         input[start:i],
			PositionIncr: 1,
			Start:        start,
			End:          i,
			Type:         analysis.AlphaNumeric,
		})

		if i-start == t.maxLength {
			start = start + 1
			i = start
		}
	}

	return rv
}
