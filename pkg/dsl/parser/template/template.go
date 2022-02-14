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

	if data["mappings"] == nil {
		return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[template] mappings should be defined")

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
			index, err := index.Request(v)
			if err != nil {
				return nil, err
			}
			template.Template = index
		case "mappings":
			v, ok := v.(map[string]interface{})
			if !ok {
				return nil, meta.NewError(meta.ErrorTypeXContentParseException, "[template] mappings value should be an object")
			}
			mappings, err := mappings.Request(v)
			if err != nil {
				return nil, err
			}
			template.Mappings = mappings
		default:
			return nil, meta.NewError(meta.ErrorTypeParsingException, fmt.Sprintf("[template] unknown option [%s]", k))
		}
	}

	return template, nil
}
