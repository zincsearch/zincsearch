package token

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"
	"golang.org/x/text/unicode/norm"

	"github.com/prabhatsharma/zinc/pkg/errors"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewUnicodenormTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	form, _ := zutils.GetStringFromMap(options, "form")
	form = strings.ToUpper(form)
	switch form {
	case "NFC":
		return token.NewUnicodeNormalizeFilter(norm.NFC), nil
	case "NFD":
		return token.NewUnicodeNormalizeFilter(norm.NFD), nil
	case "NFKC":
		return token.NewUnicodeNormalizeFilter(norm.NFKC), nil
	case "NFKD":
		return token.NewUnicodeNormalizeFilter(norm.NFKD), nil
	default:
		return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[token_filter] unicodenorm doesn't support form [%s]", form))
	}
}
