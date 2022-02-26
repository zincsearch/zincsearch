package token

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewElisionTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	articles, err := zutils.GetStringSliceFromMap(options, "articles")
	if err != nil {
		articles = []string{"l", "m", "t", "qu", "n", "s", "j", "d", "c", "jusqu", "quoiqu", "lorsqu", "puisqu"}
	}
	dict := analysis.NewTokenMap()
	for _, word := range articles {
		dict.AddToken(word)
	}
	return token.NewElisionFilter(dict), nil
}
