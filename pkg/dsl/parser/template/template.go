package template

import (
	"fmt"
	"strings"

	"github.com/prabhatsharma/zinc/pkg/dsl/meta"
	"github.com/prabhatsharma/zinc/pkg/dsl/parser/index"
	"github.com/prabhatsharma/zinc/pkg/dsl/parser/mappings"
)

func Request(data map[string]interface{}) (*meta.Template, error) {
	if data == nil {
		return nil, nil
	}

	if data["index_patterns"] == nil {
		return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[template] index_patterns should be defined")
	}

	if data["template"] == nil {
		return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[template] template should be defined")

	}

	template := new(meta.Template)
	for k, v := range data {
		k = strings.ToLower(k)
		switch k {
		case "index_patterns":
			patterns, ok := v.([]interface{})
			if !ok {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[template] index_patterns value should be an array of strings")
			}
			for _, pattern := range patterns {
				template.IndexPatterns = append(template.IndexPatterns, pattern.(string))
			}
		case "priority":
			priority, ok := v.(float64)
			if !ok {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[template] priority value should be a numberic")
			}
			template.Priority = int(priority)
		case "template":
			v, ok := v.(map[string]interface{})
			if !ok {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[template] template value should be an object")
			}
			for k, v := range v {
				k = strings.ToLower(k)
				switch k {
				case "settings":
					index, err := index.Request(map[string]interface{}{"settings": v})
					if err != nil {
						return nil, err
					}
					if index != nil {
						template.Template.Settings = index.Settings
					}
				case "mappings":
					v, ok := v.(map[string]interface{})
					if !ok {
						return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[template] mappings value should be an object")
					}
					mappings, err := mappings.Request(v)
					if err != nil {
						return nil, err
					}
					if mappings != nil {
						template.Template.Mappings = mappings
					}
				default:
					return nil, meta.NewError(meta.ErrorTypeXContentParseException, fmt.Sprintf("[template] template unknown option [%s]", k))
				}
			}
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[template] unknown option [%s]", k))
		}
	}

	return template, nil
}
