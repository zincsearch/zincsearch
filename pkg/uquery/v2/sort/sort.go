package sort

import (
	"strings"

	"github.com/blugelabs/bluge/search"
	"github.com/prabhatsharma/zinc/pkg/errors"
)

func Request(v interface{}) (search.SortOrder, error) {
	if v == nil {
		return nil, nil
	}

	sorts := make(search.SortOrder, 0, 1)
	switch v := v.(type) {
	case string:
		sorts = append(sorts, search.ParseSearchSortString(v))
		return sorts, nil
	case []interface{}:
		for _, v := range v {
			switch v := v.(type) {
			case string:
				sorts = append(sorts, search.ParseSearchSortString(v))
			case map[string]interface{}:
				if len(v) > 1 {
					return nil, errors.New(errors.ErrorTypeParsingException, "[sort] field doesn't support multiple values")
				}
				for field, v := range v {
					sort := search.SortBy(search.Field(field))
					switch v := v.(type) {
					case string:
						if strings.ToLower(v) == "desc" {
							sort.Desc()
						}
					case map[string]interface{}:
						for kk, vv := range v {
							kk = strings.ToLower(kk)
							switch kk {
							case "order":
								if strings.ToLower(vv.(string)) == "desc" {
									sort.Desc()
								}
							case "format":
							default:
							}
						}
					default:
					}
					sorts = append(sorts, sort)
				}
			}
		}
	default:
		return nil, errors.New(errors.ErrorTypeXContentParseException, "[sort] value should be string or array")
	}

	return sorts, nil
}
