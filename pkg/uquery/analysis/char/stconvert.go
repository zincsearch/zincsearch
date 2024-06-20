package char

import (
	"github.com/blugelabs/bluge/analysis"
	zincchar "github.com/zincsearch/zincsearch/pkg/bluge/analysis/char"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func NewSTConvertCharFilter(options interface{}) (analysis.CharFilter, error) {
	conversion, err := zutils.GetStringFromMap(options, "convert_type")
	if err != nil {
		return nil, err
	}
	return zincchar.NewSTConvertCharFilter(conversion)
}
