package token

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/go-ego/gse"
)

type StopTokenFilter struct {
	seg *gse.Segmenter
}

func NewStopTokenFilter(seg *gse.Segmenter, stopwords []string) *StopTokenFilter {
	if len(stopwords) > 0 {
		for _, word := range stopwords {
			seg.AddStop(word)
		}
	}
	return &StopTokenFilter{seg}
}

func (f *StopTokenFilter) Filter(input analysis.TokenStream) analysis.TokenStream {
	var j, skipped int
	for _, token := range input {
		if !f.seg.IsStop(string(token.Term)) {
			token.PositionIncr += skipped
			skipped = 0
			input[j] = token
			j++
		} else {
			skipped += token.PositionIncr
		}
	}

	return input[:j]
}
