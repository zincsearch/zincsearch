package tokenizer

import (
	"github.com/blugelabs/bluge/analysis"

	zinctokenizer "github.com/prabhatsharma/zinc/pkg/bluge/analysis/tokenizer"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func NewPathHierarchyTokenizer(options interface{}) (analysis.Tokenizer, error) {
	delimiter, _ := zutils.GetStringFromMap(options, "delimiter")
	if len(delimiter) == 0 {
		delimiter = "/"
	}
	replacement, _ := zutils.GetStringFromMap(options, "replacement")
	if len(replacement) == 0 {
		replacement = delimiter
	}
	skip, _ := zutils.GetFloatFromMap(options, "skip")
	return zinctokenizer.NewPathHierarchyTokenizer(delimiter[0], replacement[0], int(skip)), nil
}
