package highlight

import (
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func Request(highlight *meta.Highlight) error {
	if len(highlight.Fields) == 0 {
		return nil
	}

	if highlight.NumberOfFragments == 0 {
		highlight.NumberOfFragments = 3
	}
	for _, field := range highlight.Fields {
		if field.FragmentSize == 0 && highlight.FragmentSize > 0 {
			field.FragmentSize = highlight.FragmentSize
		}
		if field.NumberOfFragments == 0 && highlight.NumberOfFragments > 0 {
			field.NumberOfFragments = highlight.NumberOfFragments
		}
	}

	return nil
}
