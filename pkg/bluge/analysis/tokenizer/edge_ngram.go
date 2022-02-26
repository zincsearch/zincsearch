package tokenizer

import (
	"github.com/blugelabs/bluge/analysis"
)

type EdgeNgramTokenizer struct {
	minLength  int
	maxLength  int
	tokenChars []string
}

func NewEdgeNgramTokenizer(minLength, maxLength int, tokenChars []string) *EdgeNgramTokenizer {
	return &EdgeNgramTokenizer{
		minLength:  minLength,
		maxLength:  maxLength,
		tokenChars: tokenChars,
	}
}

func (t *EdgeNgramTokenizer) Tokenize(input []byte) analysis.TokenStream {
	n := len(input)
	start := 0
	rv := make(analysis.TokenStream, 0, n)
	for i := 1; i <= n; i++ {
		if i-start >= t.minLength {
			valid := true
			if len(t.tokenChars) > 0 {
				for _, c := range string(input[start:i]) {
					if !t.isChar(c) {
						valid = false
						break
					}
				}
			}
			if valid {
				rv = append(rv, &analysis.Token{
					Term:         input[start:i],
					PositionIncr: 1,
					Start:        start,
					End:          i,
					Type:         analysis.AlphaNumeric,
				})
			} else {
				start = i
				continue
			}
		}

		if i-start == t.maxLength {
			break
		}
	}

	return rv
}

func (t *EdgeNgramTokenizer) isChar(r rune) bool {
	var ok bool
	for _, char := range t.tokenChars {
		if ok = isChar(char, r); ok {
			return true
		}
	}

	return false
}
