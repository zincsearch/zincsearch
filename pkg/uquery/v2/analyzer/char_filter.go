package analyzer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/char"

	zincchar "github.com/prabhatsharma/zinc/pkg/bluge/analysis/char"
	"github.com/prabhatsharma/zinc/pkg/errors"
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
			pattern, err := zutils.GetStringFromMap(options, "pattern")
			if err != nil {
				return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[char_filter] %s option [%s] should be exists", filterType, "pattern"))
			}
			replacement, err := zutils.GetStringFromMap(options, "replacement")
			if err != nil {
				return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[char_filter] %s option [%s] should be exists", filterType, "replacement"))
			}
			re := regexp.MustCompile(pattern)
			filters[name] = char.NewRegexpCharFilter(re, []byte(replacement))
		case "mapping":
			mappings, err := zutils.GetStringSliceFromMap(options, "mappings")
			if err != nil {
				return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[char_filter] %s option [%s] should be exists", filterType, "mappings"))
			}
			for _, mapping := range mappings {
				if !strings.Contains(mapping, " => ") {
					return nil, errors.New(errors.ErrorTypeRuntimeException, fmt.Sprintf("[char_filter] %s option [%s] Invalid Mapping Rule: [%s]", filterType, "mappings", mapping))
				}
			}
			filters[name] = zincchar.NewMappingCharFilter(mappings)
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[char_filter] doesn't support type [%s]", filterType))
		}
	}

	return filters, nil
}
