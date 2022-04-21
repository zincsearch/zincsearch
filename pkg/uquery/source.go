package uquery

import (
	"github.com/goccy/go-json"

	v1 "github.com/zinclabs/zinc/pkg/meta/v1"
)

func HandleSource(source *v1.Source, data []byte) map[string]interface{} {
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

	// delete field not in source.Fields
	for field := range ret {
		if _, ok := source.Fields[field]; ok {
			continue
		}
		delete(ret, field)
	}

	return ret
}
