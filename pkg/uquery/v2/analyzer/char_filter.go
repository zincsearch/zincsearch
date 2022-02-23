package analyzer

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/char"

	"github.com/prabhatsharma/zinc/pkg/errors"
	zincchar "github.com/prabhatsharma/zinc/pkg/uquery/v2/analyzer/char"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func RequestCharFilter(data map[string]interface{}) (map[string]analysis.CharFilter, error) {
	if data == nil {
		return nil, nil
	}

	filters := make(map[string]analysis.CharFilter)
	for name, options := range data {
		typ, err := zutils.GetStringFromMap(options, "type")
		if err != nil {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[char_filter] %s option [%s] should be exists", name, "type"))
		}
		filter, err := RequestCharFilterSingle(typ, options)
		if err != nil {
			return nil, err
		}
		filters[name] = filter
	}

	return filters, nil
}

func RequestCharFilterSlice(data []interface{}) ([]analysis.CharFilter, error) {
	if data == nil {
		return nil, nil
	}

	filters := make([]analysis.CharFilter, 0, len(data))
	for _, typ := range data {
		typ, ok := typ.(string)
		if !ok {
			return nil, errors.New(errors.ErrorTypeParsingException, "[char_filter] option should be string")
		}
		filter, err := RequestCharFilterSingle(typ, nil)
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}

	return filters, nil
}

func RequestCharFilterSingle(typ string, options interface{}) (analysis.CharFilter, error) {
	typ = strings.ToLower(typ)
	switch typ {
	case "ascii_folding":
		return char.NewASCIIFoldingFilter(), nil
	case "html", "html_strip":
		return char.NewHTMLCharFilter(), nil
	case "zero_width_non_joiner":
		return char.NewZeroWidthNonJoinerCharFilter(), nil
	case "regexp", "pattern_replace":
		return zincchar.NewRegexpCharFilter(options)
	case "mapping":
		return zincchar.NewMappingCharFilter(options)
	default:
		return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[char_filter] doesn't support filter [%s]", typ))
	}
}
