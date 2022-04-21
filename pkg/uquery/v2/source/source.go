package source

import (
	"strings"

	"github.com/goccy/go-json"

	"github.com/zinclabs/zinc/pkg/errors"
	meta "github.com/zinclabs/zinc/pkg/meta/v2"
)

func Request(v interface{}) (*meta.Source, error) {
	source := &meta.Source{Enable: true}
	if v == nil {
		return source, nil
	}

	switch v := v.(type) {
	case bool:
		source.Enable = v
	case []interface{}:
		source.Fields = make([]string, 0, len(v))
		for _, field := range v {
			if v, ok := field.(string); ok {
				source.Fields = append(source.Fields, v)
			} else {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[_source] value should be boolean or []string")
			}
		}
	default:
		return nil, errors.New(errors.ErrorTypeXContentParseException, "[_source] value should be boolean or []string")
	}

	return source, nil
}

func Response(source *meta.Source, data []byte) map[string]interface{} {
	// return empty
	if !source.Enable {
		return nil
	}

	ret := make(map[string]interface{})
	err := json.Unmarshal(data, &ret)
	if err != nil {
		return nil
	}

	// return all fields
	if len(source.Fields) == 0 {
		return ret
	}

	wildcard := false
	rets := make(map[string]interface{})
	for _, field := range source.Fields {
		wildcard = false
		if strings.HasSuffix(field, "*") {
			wildcard = true
		}
		if _, ok := ret[field]; ok {
			rets[field] = ret[field]
		} else if wildcard {
			for k, v := range ret {
				if strings.HasPrefix(k, field[:len(field)-1]) {
					rets[k] = v
				}
			}
		}
	}

	return rets
}
