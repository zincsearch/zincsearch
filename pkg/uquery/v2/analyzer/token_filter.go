package analyzer

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/prabhatsharma/zinc/pkg/errors"
	zinctoken "github.com/prabhatsharma/zinc/pkg/uquery/v2/analyzer/token"
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
			filters[name] = token.NewCamelCaseFilter()
		case "dict":
			filters[name], err = zinctoken.NewDictTokenFilter(options)
		case "edge_ngram":
			filters[name], err = zinctoken.NewEdgeNgramTokenFilter(options)
		case "elision":
			filters[name], err = zinctoken.NewElisionTokenFilter(options)
		case "keyword":
			filters[name], err = zinctoken.NewKeywordTokenFilter(options)
		case "length":
			filters[name], err = zinctoken.NewLengthTokenFilter(options)
		case "lower_case":
			filters[name] = token.NewLowerCaseFilter()
		case "ngram":
			filters[name], err = zinctoken.NewNgramTokenFilter(options)
		case "porter":
			filters[name] = token.NewPorterStemmer()
		case "reverse":
			filters[name] = token.NewReverseFilter()
		case "shingle":
			filters[name], err = zinctoken.NewShingleTokenFilter(options)
		case "stop":
			filters[name], err = zinctoken.NewStopTokenFilter(options)
		case "truncate":
			filters[name], err = zinctoken.NewTruncateTokenFilter(options)
		case "unicodenorm":
			filters[name], err = zinctoken.NewUnicodenormTokenFilter(options)
		case "unique":
			filters[name] = token.NewUniqueTermFilter()
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[token_filter] doesn't support type [%s]", filterType))
		}

		if err != nil {
			return nil, err
		}
	}

	return filters, nil
}
