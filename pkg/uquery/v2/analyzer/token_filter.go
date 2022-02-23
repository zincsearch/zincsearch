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
		typ, err := zutils.GetStringFromMap(options, "type")
		if err != nil {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[token_filter] %s option [%s] should be exists", name, "type"))
		}
		filter, err := RequestTokenFilterSingle(typ, options)
		if err != nil {
			return nil, err
		}
		filters[name] = filter
	}

	return filters, nil
}

func RequestTokenFilterSlice(data []interface{}) ([]analysis.TokenFilter, error) {
	if data == nil {
		return nil, nil
	}

	filters := make([]analysis.TokenFilter, 0, len(data))
	for _, typ := range data {
		typ, ok := typ.(string)
		if !ok {
			return nil, errors.New(errors.ErrorTypeParsingException, "[token_filter] option should be string")
		}
		filter, err := RequestTokenFilterSingle(typ, nil)
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}

	return filters, nil
}

func RequestTokenFilterSingle(typ string, options interface{}) (analysis.TokenFilter, error) {
	typ = strings.ToLower(typ)
	switch typ {
	case "apostrophe":
		return token.NewApostropheFilter(), nil
	case "camel_case":
		return token.NewCamelCaseFilter(), nil
	case "dict":
		return zinctoken.NewDictTokenFilter(options)
	case "edge_ngram":
		return zinctoken.NewEdgeNgramTokenFilter(options)
	case "elision":
		return zinctoken.NewElisionTokenFilter(options)
	case "keyword":
		return zinctoken.NewKeywordTokenFilter(options)
	case "length":
		return zinctoken.NewLengthTokenFilter(options)
	case "lower_case":
		return token.NewLowerCaseFilter(), nil
	case "ngram":
		return zinctoken.NewNgramTokenFilter(options)
	case "porter":
		return token.NewPorterStemmer(), nil
	case "reverse":
		return token.NewReverseFilter(), nil
	case "shingle":
		return zinctoken.NewShingleTokenFilter(options)
	case "stop":
		return zinctoken.NewStopTokenFilter(options)
	case "truncate":
		return zinctoken.NewTruncateTokenFilter(options)
	case "unicodenorm":
		return zinctoken.NewUnicodenormTokenFilter(options)
	case "unique":
		return token.NewUniqueTermFilter(), nil
	default:
		return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[token_filter] doesn't support filter [%s]", typ))
	}
}
