package tokenizer

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"
	"github.com/blugelabs/bluge/analysis/tokenizer"
)

type LowerCaseTokenizer struct{}

func NewLowerCaseTokenizer() *LowerCaseTokenizer {
	return &LowerCaseTokenizer{}

}

func (t *LowerCaseTokenizer) Tokenize(input []byte) analysis.TokenStream {
	tokens := tokenizer.NewLetterTokenizer().Tokenize(input)
	filter := token.NewLowerCaseFilter()
	filter.Filter(tokens)
	return tokens
}
