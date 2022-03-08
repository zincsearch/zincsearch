package tokenizer

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/go-ego/gse"
)

type SearchTokenizer struct {
	seg *gse.Segmenter
}

func NewSearchTokenizer(seg *gse.Segmenter) *SearchTokenizer {
	return &SearchTokenizer{seg}
}

func (t *SearchTokenizer) Tokenize(input []byte) analysis.TokenStream {
	result := make(analysis.TokenStream, 0, len(input))
	text := string(input)
	search := t.seg.CutSearch(text, true)
	tokens := t.seg.Analyze(search, text)
	var start, positionIncr int
	for _, token := range tokens {
		positionIncr = 1
		if start == token.Start {
			positionIncr = 0
		}
		start = token.Start

		typ := analysis.Ideographic
		alphaNumeric := true
		for _, r := range token.Text {
			if r < 32 || r > 126 {
				alphaNumeric = false
				break
			}
		}
		if alphaNumeric {
			typ = analysis.AlphaNumeric
		}

		result = append(result, &analysis.Token{
			Term:         []byte(token.Text),
			Start:        token.Start,
			End:          token.End,
			PositionIncr: positionIncr,
			Type:         typ,
		})
	}
	return result
}
