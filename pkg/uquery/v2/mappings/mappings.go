package mappings

import (
	"fmt"
	"strings"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func Request(data map[string]interface{}) (*meta.Mappings, error) {
	if data == nil {
		return nil, nil
	}

	if data["properties"] == nil {
		return nil, errors.New(errors.ErrorTypeParsingException, "[mappings] properties should be defined")

	}

	properties, ok := data["properties"].(map[string]interface{})
	if !ok {
		return nil, errors.New(errors.ErrorTypeParsingException, "[mappings] properties should be an object")
	}

	mappings := meta.NewMappings()
	for field, prop := range properties {
		prop, ok := prop.(map[string]interface{})
		if !ok {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[mappings] properties [%s] should be an object", field))
		}
		propType, ok := prop["type"]
		if !ok {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[mappings] properties [%s] should be exists", "type"))
		}
		propTypeStr, ok := propType.(string)
		if !ok {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[mappings] properties [%s] should be an string", "type"))
		}

		var newProp meta.Property
		propTypeStr = strings.ToLower(propTypeStr)
		switch propTypeStr {
		case "text", "keyword", "numeric", "bool", "time":
			newProp = meta.NewProperty(propTypeStr)
		case "integer", "double", "long":
			newProp = meta.NewProperty("numeric")
		case "boolean":
			newProp = meta.NewProperty("bool")
		case "date", "datetime":
			newProp = meta.NewProperty("time")
		case "flattened", "object", "match_only_text":
			// ignore
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[mappings] properties [%s] doesn't support type [%s]", field, propTypeStr))
		}

		for k, v := range prop {
			switch k {
			case "type":
				// handled
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
				// ignore unknown options
				// return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[mappings] properties [%s] unknown option [%s]", field, k))
			}
		}

		if newProp.Type != "" {
			mappings.Properties[field] = newProp
		}
	}

	return mappings, nil
}
