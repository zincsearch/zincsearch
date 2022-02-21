package analyzer

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"
	"github.com/prabhatsharma/zinc/pkg/errors"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func RequestTokenFilter(data map[string]interface{}) (map[string]analysis.TokenFilter, error) {
	if data == nil {
		return nil, nil
	}

	filters := make(map[string]analysis.TokenFilter)
	for name, options := range data {
		filterType, err := zutils.GetStringFromMap(options, "type")
		if err != nil {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[token_filter] %s option [%s] should be exists", name, "type"))
		}
		filterType = strings.ToLower(filterType)
		switch filterType {
		case "apostrophe":
			filters[name] = token.NewApostropheFilter()
		case "camel_case":
		case "dict":
		case "edge_ngram":
		case "elision":
		case "keyword":
		case "length":
		case "lower_case":
		case "ngram":
		case "porter":
		case "reverse":
		case "shingle":
		case "stop":
		case "truncate":
		case "unicodenorm":
		case "unique":
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[token_filter] doesn't support type [%s]", filterType))
		}
	}

	return filters, nil
}
