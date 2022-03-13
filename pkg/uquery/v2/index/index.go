package index

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/analysis"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/mappings"
)

func Request(data map[string]interface{}) (*meta.Index, error) {
	if len(data) == 0 {
		return nil, nil
	}

	index := new(meta.Index)
	for k, v := range data {
		k = strings.ToLower(k)
		switch k {
		case "settings":
			v, ok := v.(map[string]interface{})
			if !ok {
				return nil, errors.New(errors.ErrorTypeParsingException, "[index] settings should be an object")
			}
			vjson, _ := json.Marshal(v)
			settings := new(meta.IndexSettings)
			if err := json.Unmarshal(vjson, settings); err != nil {
				return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[index] settings parse error: %s", err.Error()))
			}
			if _, err := analysis.RequestAnalyzer(settings.Analysis); err != nil {
				return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[index] settings.analysis parse error: %s", err.Error()))
			}
			if settings != nil && (settings.NumberOfShards > 0 || settings.NumberOfReplicas > 0 || settings.Analysis != nil) {
				index.Settings = settings
			}
		case "mappings":
			v, ok := v.(map[string]interface{})
			if !ok {
				return nil, errors.New(errors.ErrorTypeParsingException, "[index] mappings should be an object")
			}
			mappings, err := mappings.Request(v)
			if err != nil {
				return nil, err
			}
			index.Mappings = mappings
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[index] unknown option [%s]", k))
		}
	}

	return index, nil
}
