package tokenizer

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/go-ego/gse"
)

type StandardTokenizer struct {
	seg *gse.Segmenter
}

func NewStandardTokenizer(seg *gse.Segmenter) *StandardTokenizer {
	return &StandardTokenizer{seg}
}

func (t *StandardTokenizer) Tokenize(input []byte) analysis.TokenStream {
	result := make(analysis.TokenStream, 0, len(input))
	segments := t.seg.Segment(input)
	for _, seg := range segments {
		typ := analysis.Ideographic
		alphaNumeric := true
		for _, r := range seg.Token().Text() {
			if r < 32 || r > 126 {
				alphaNumeric = false
				break
			}
		}
		if alphaNumeric {
			typ = analysis.AlphaNumeric
		}
		result = append(result, &analysis.Token{
			Term:         []byte(seg.Token().Text()),
			Start:        seg.Start(),
			End:          seg.End(),
			PositionIncr: 1,
			Type:         typ,
		})
	}
	return result
}
