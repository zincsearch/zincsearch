package mappings

import (
	"fmt"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func Request(data map[string]interface{}) (*meta.Mappings, error) {
	if data == nil {
		return nil, nil
	}

	if _, ok := data["properties"]; !ok {
		return nil, nil
	}

	properties, ok := data["properties"].(map[string]interface{})
	if !ok {
		return nil, meta.NewError(meta.ErrorTypeParsingException, "mappings.properties should be a object")
	}

	mappings := new(meta.Mappings)
	mappings.Properties = make(map[string]meta.Property)
	for field, prop := range properties {
		prop, ok := prop.(map[string]interface{})
		if !ok {
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("mappings.properties [%s] should be a object", field))
		}
		newProp := meta.NewProperty("text")
		for k, v := range prop {
			switch k {
			case "type":
				newProp.Type = v.(string)
			case "analyzer":
				newProp.Analyzer = v.(string)
			case "search_analyzer":
				newProp.SearchAnalyzer = v.(string)
			case "format":
				newProp.Format = v.(string)
			case "index":
				newProp.Index = v.(bool)
			case "store":
				newProp.Store = v.(bool)
			case "sortable":
				newProp.Sortable = v.(bool)
			case "aggregatable":
				newProp.Aggregatable = v.(bool)
			case "highlightable":
				newProp.Highlightable = v.(bool)
			default:
				return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("mappings.properties [%s] unknown option [%s]", field, k))
			}
		}
		switch newProp.Type {
		case "text", "keyword", "numeric", "bool", "time":
			// continue
		case "integer", "double", "long":
			newProp.Type = "numeric"
		case "boolean":
			newProp.Type = "bool"
		case "date", "datetime":
			newProp.Type = "time"
		default:
			return nil, fmt.Errorf("mappings.properties [%s] doesn't support type [%s]", newProp.Type, field)
		}
		mappings.Properties[field] = newProp
	}

	return mappings, nil
}
