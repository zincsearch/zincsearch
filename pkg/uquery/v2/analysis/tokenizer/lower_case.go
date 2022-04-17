package tokenizer

import (
	"github.com/blugelabs/bluge/analysis"

	zinctokenizer "github.com/zinclabs/zinc/pkg/bluge/analysis/tokenizer"
)

func NewLowerCaseTokenizer() (analysis.Tokenizer, error) {
	return zinctokenizer.NewLowerCaseTokenizer(), nil
}
