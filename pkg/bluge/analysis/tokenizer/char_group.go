package tokenizer

import (
	"unicode"

	"github.com/blugelabs/bluge/analysis/tokenizer"
)

type IsTokenRune func(r rune) bool

type CharGroupTokenizer struct {
	chars []string
}

func NewCharGroupTokenizer(chars []string) *tokenizer.CharacterTokenizer {
	t := new(CharGroupTokenizer)
	for _, char := range chars {
		if len(char) == 0 {
			continue
		}
		t.chars = append(t.chars, char)
	}
	return tokenizer.NewCharacterTokenizer(t.isChar)
}

func (t *CharGroupTokenizer) isChar(r rune) bool {
	var ok bool
	for _, char := range t.chars {
		if ok = isChar(char, r); ok {
			return false
		}
	}

	return true
}

func isChar(char string, r rune) bool {
	ok := false
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

	return ok
}
