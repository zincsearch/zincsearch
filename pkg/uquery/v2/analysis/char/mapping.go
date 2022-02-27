package char

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge/analysis"

	zincchar "github.com/prabhatsharma/zinc/pkg/bluge/analysis/char"
	"github.com/prabhatsharma/zinc/pkg/errors"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewMappingCharFilter(options interface{}) (analysis.CharFilter, error) {
	mappings, err := zutils.GetStringSliceFromMap(options, "mappings")
	if err != nil || len(mappings) == 0 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[char_filter] mapping option [mappings] should be exists")
	}
	for _, mapping := range mappings {
		if !strings.Contains(mapping, " => ") {
			return nil, errors.New(errors.ErrorTypeRuntimeException, fmt.Sprintf("[char_filter] mapping option [mappings] Invalid Mapping Rule: [%s], should be [old => new]", mapping))
		}
	}

	return zincchar.NewMappingCharFilter(mappings), nil
}
