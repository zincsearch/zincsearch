package fields

import (
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/zinclabs/zinc/pkg/errors"
	meta "github.com/zinclabs/zinc/pkg/meta/v2"
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
					return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[fields] unknown field [%s]", k))
				}
			}
			fields = append(fields, f)
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, "[fields] value should be string or object")

		}
	}

	return fields, nil
}

func Response(fields []*meta.Field, data []byte, mappings *meta.Mappings) map[string]interface{} {
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
	results := make(map[string]interface{})
	for _, v := range fields {
		wildcard = false
		field = v.Field
		if strings.HasSuffix(field, "*") {
			wildcard = true
		}
		if rv, ok := ret[field]; ok {
			if mappings.Properties[field].Type == "time" && v.Format != "" {
				if t, err := time.Parse(mappings.Properties[field].Format, rv.(string)); err == nil {
					results[field] = []interface{}{t.Format(v.Format)}
				} else {
					results[field] = []interface{}{rv}
				}
			} else {
				results[field] = []interface{}{rv}
			}
		} else if wildcard {
			for rk, rv := range ret {
				if strings.HasPrefix(rk, field[:len(field)-1]) {
					if mappings.Properties[rk].Type == "time" && v.Format != "" {
						if t, err := time.Parse(mappings.Properties[rk].Format, rv.(string)); err == nil {
							results[rk] = []interface{}{t.Format(v.Format)}
						} else {
							results[rk] = []interface{}{rv}
						}
					} else {
						results[rk] = []interface{}{rv}
					}
				}
			}
		}
	}

	return results
}
