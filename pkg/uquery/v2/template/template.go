package template

import (
	"fmt"
	"strings"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/index"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/mappings"
)

func Request(data map[string]interface{}) (*meta.Template, error) {
	if data == nil {
		return nil, nil
	}

	if data["template"] == nil {
		return nil, errors.New(errors.ErrorTypeXContentParseException, "[template] template should be defined")
	}

	template := new(meta.Template)
	for k, v := range data {
		k = strings.ToLower(k)
		switch k {
		case "index_patterns":
			patterns, ok := v.([]interface{})
			if !ok {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[template] index_patterns value should be an array of string")
			}
			for _, pattern := range patterns {
				template.IndexPatterns = append(template.IndexPatterns, pattern.(string))
			}
		case "priority":
			priority, ok := v.(float64)
			if !ok {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[template] priority value should be a numberic")
			}
			template.Priority = int(priority)
		case "template":
			switch v := v.(type) {
			case string:
				// compatible {"priority":150,"template":"filebeat-7.16.3-*"}
				template.IndexPatterns = append(template.IndexPatterns, v)
			case map[string]interface{}:
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
							return nil, errors.New(errors.ErrorTypeXContentParseException, "[template] mappings value should be an object")
						}
						mappings, err := mappings.Request(v)
						if err != nil {
							return nil, err
						}
						if mappings != nil {
							template.Template.Mappings = mappings
						}
					case "alias":
						// TODO: implement
					default:
						return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[template] template unknown option [%s]", k))
					}
				}
			default:
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[template] template value should be an object")
			}
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[template] unknown option [%s]", k))
		}
	}

	if len(template.IndexPatterns) == 0 {
		return nil, errors.New(errors.ErrorTypeXContentParseException, "[template] index_patterns should be defined")
	}

	return template, nil
}
