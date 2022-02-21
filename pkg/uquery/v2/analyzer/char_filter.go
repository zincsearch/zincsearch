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
		filterType, err := zutils.GetStringFromMap(options, "type")
		if err != nil {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[char_filter] %s option [%s] should be exists", name, "type"))
		}
		filterType = strings.ToLower(filterType)
		switch filterType {
		case "ascii_folding":
			filters[name] = char.NewASCIIFoldingFilter()
		case "html", "html_strip":
			filters[name] = char.NewHTMLCharFilter()
		case "zero_width_non_joiner":
			filters[name] = char.NewZeroWidthNonJoinerCharFilter()
		case "regexp", "pattern_replace":
			filters[name], err = zincchar.NewRegexpCharFilter(options)
		case "mapping":
			filters[name], err = zincchar.NewMappingCharFilter(options)
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[char_filter] doesn't support filter [%s]", filterType))
		}

		if err != nil {
			return nil, err
		}
	}

	return filters, nil
}

func RequestCharFilterSlice(data []interface{}) ([]analysis.CharFilter, error) {
	if data == nil {
		return nil, nil
	}

	filters := make([]analysis.CharFilter, 0, len(data))
	for _, name := range data {
		name, ok := name.(string)
		if !ok {
			return nil, errors.New(errors.ErrorTypeParsingException, "[char_filter] option should be string")
		}
		name = strings.ToLower(name)
		switch name {
		case "ascii_folding":
			filters = append(filters, char.NewASCIIFoldingFilter())
		case "html", "html_strip":
			filters = append(filters, char.NewHTMLCharFilter())
		case "zero_width_non_joiner":
			filters = append(filters, char.NewZeroWidthNonJoinerCharFilter())
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[char_filter] doesn't support filter [%s]", name))
		}
	}

	return filters, nil
}
