package tokenizer

import (
	"github.com/blugelabs/bluge/analysis"
)

type PathHierarchyTokenizer struct {
	dilimiter   byte
	replacement byte
	skip        int
}

func NewPathHierarchyTokenizer(dilimiter, replacement byte, skip int) *PathHierarchyTokenizer {
	if dilimiter == 0 {
		dilimiter = '/'
	}
	if replacement == 0 {
		replacement = dilimiter
	}
	return &PathHierarchyTokenizer{
		dilimiter:   dilimiter,
		replacement: replacement,
		skip:        skip,
	}
}

func (t *PathHierarchyTokenizer) Tokenize(input []byte) analysis.TokenStream {
	n := len(input)
	rv := make(analysis.TokenStream, 0, n)
	start := 0
	skip := 0
	for i := 1; i < n; i++ {
		if input[i] == t.dilimiter {
			input[i] = t.replacement
			if t.skip > 0 && skip < t.skip {
				skip++
				start = i
				continue
			}
			rv = append(rv, t.makeToken(input, start, i))
		}
	}

	if input[n-1] != t.dilimiter && skip == t.skip {
		rv = append(rv, t.makeToken(input, start, n))
	}

	return rv
}

func (t *PathHierarchyTokenizer) makeToken(input []byte, start, end int) *analysis.Token {
	return &analysis.Token{
		Term:         input[start:end],
		PositionIncr: 1,
		Start:        start,
		End:          end,
		Type:         analysis.AlphaNumeric,
	}
}
