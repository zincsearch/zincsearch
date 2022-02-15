package fields

import (
	"encoding/json"
	"fmt"
	"strings"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func Request(v []interface{}) ([]*meta.Field, error) {
	if v == nil {
		return nil, nil
	}

	fields := make([]*meta.Field, 0, len(v))
	for _, v := range v {
		switch v := v.(type) {
		case string:
			fields = append(fields, &meta.Field{Field: v})
		case map[string]interface{}:
			f := new(meta.Field)
			for k, v := range v {
				k = strings.ToLower(k)
				switch k {
				case "field":
					f.Field = v.(string)
				case "format":
					f.Format = v.(string)
				default:
					return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[fields] unknown field [%s]", k))
				}
			}
			fields = append(fields, f)
		default:
			return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[fields] value should be string or object")

		}
	}

	return fields, nil
}

func Response(fields []*meta.Field, data []byte) map[string]interface{} {
	// return empty
	if len(fields) == 0 {
		return nil
	}

	ret := make(map[string]interface{})
	err := json.Unmarshal(data, &ret)
	if err != nil {
		return nil
	}

	var field string
	wildcard := false
	rets := make(map[string]interface{})
	for _, v := range fields {
		wildcard = false
		field = v.Field
		if strings.HasSuffix(field, "*") {
			wildcard = true
		}
		if _, ok := ret[field]; ok {
			rets[field] = []interface{}{ret[field]}
		} else if wildcard {
			for k, v := range ret {
				if strings.HasPrefix(k, field[:len(field)-1]) {
					rets[k] = []interface{}{v}
				}
			}
		}
	}

	// TODO: field format

	return rets
}
