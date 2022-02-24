package tokenizer

import (
	"unicode"

	"github.com/blugelabs/bluge/analysis/tokenizer"
)

type IsTokenRune func(r rune) bool

type CharGroupTokenizer struct {
	cahrs []string
}

func NewCharGroupTokenizer(chars []string) *tokenizer.CharacterTokenizer {
	cg := new(CharGroupTokenizer)
	for _, char := range chars {
		if len(char) == 0 {
			continue
		}
		cg.cahrs = append(cg.cahrs, char)
	}
	return tokenizer.NewCharacterTokenizer(cg.isChar)
}

func (cg *CharGroupTokenizer) isChar(r rune) bool {
	var ok bool
	for _, char := range cg.cahrs {
		switch char {
		case "graphic":
			ok = unicode.IsGraphic(r)
		case "print":
			ok = unicode.IsPrint(r)
		case "control":
			ok = unicode.IsControl(r)
		case "letter":
			ok = unicode.IsLetter(r)
		case "mark":
			ok = unicode.IsMark(r)
		case "number", "digit":
			ok = unicode.IsNumber(r)
		case "punct", "punctuation":
			ok = unicode.IsPunct(r)
		case "space", "whitespace", "white_space":
			ok = unicode.IsSpace(r)
		case "symbol":
			ok = unicode.IsSymbol(r)
		default:
			for _, c := range char {
				if r == c {
					ok = true
					break
				}
			}
		}
		if ok {
			return false
		}
	}

	return true
}
