package source

import (
	"encoding/json"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func Request(s interface{}) (*meta.Source, error) {
	source := &meta.Source{Enable: true}
	switch v := s.(type) {
	case bool:
		source.Enable = v
	case []interface{}:
		source.Fields = make(map[string]bool, len(v))
		for _, field := range v {
			if fv, ok := field.(string); ok {
				source.Fields[fv] = true
			} else {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[_source] value should be boolean or []string")
			}
		}
	default:
		return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[_source] value should be boolean or []string")
	}

	return source, nil
}

func Response(s interface{}, data []byte) map[string]interface{} {
	source, ok := s.(*meta.Source)
	if !ok {
		source = &meta.Source{Enable: true}
	}

	ret := make(map[string]interface{})
	// return empty
	if !source.Enable {
		return ret
	}

	err := json.Unmarshal(data, &ret)
	if err != nil {
		return nil
	}

	// return all fields
	if len(source.Fields) == 0 {
		return ret
	}

	// TODO: wildcard support
	// delete field not in source.Fields
	for field := range ret {
		if _, ok := source.Fields[field]; ok {
			continue
		}
		delete(ret, field)
	}

	return ret
}
