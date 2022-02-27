package tokenizer

import (
	"fmt"
	"unicode"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/tokenizer"

	"github.com/prabhatsharma/zinc/pkg/errors"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewCharacterTokenizer(options interface{}) (analysis.Tokenizer, error) {
	char, _ := zutils.GetStringFromMap(options, "char")
	switch char {
	case "graphic":
		return tokenizer.NewCharacterTokenizer(unicode.IsGraphic), nil
	case "print":
		return tokenizer.NewCharacterTokenizer(unicode.IsPrint), nil
	case "control":
		return tokenizer.NewCharacterTokenizer(unicode.IsControl), nil
	case "letter":
		return tokenizer.NewCharacterTokenizer(unicode.IsLetter), nil
	case "mark":
		return tokenizer.NewCharacterTokenizer(unicode.IsMark), nil
	case "number", "digit":
		return tokenizer.NewCharacterTokenizer(unicode.IsNumber), nil
	case "punct", "punctuation":
		return tokenizer.NewCharacterTokenizer(unicode.IsPunct), nil
	case "space", "whitespace", "white_space":
		return tokenizer.NewCharacterTokenizer(unicode.IsSpace), nil
	case "symbol":
		return tokenizer.NewCharacterTokenizer(unicode.IsSymbol), nil
	default:
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[tokenizer] character doesn't support char [%s]", char))
	}
}
