package core

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/jeremywohl/flatten"
	"github.com/rs/zerolog/log"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	zincanalysis "github.com/prabhatsharma/zinc/pkg/uquery/v2/analysis"
)

// BuildBlugeDocumentFromJSON returns the bluge document for the json document. It also updates the mapping for the fields if not found.
// If no mappings are found, it creates te mapping for all the encountered fields. If mapping for some fields is found but not for others
// then it creates the mapping for the missing fields.
func (index *Index) BuildBlugeDocumentFromJSON(docID string, doc *map[string]interface{}) (*bluge.Document, error) {
	// Pick the index mapping from the cache if it already exists
	mappings := index.CachedMappings
	if mappings == nil {
		mappings = meta.NewMappings()
	}

	mappingsNeedsUpdate := false

	// Create a new bluge document
	bdoc := bluge.NewDocument(docID)
	flatDoc, _ := flatten.Flatten(*doc, "", flatten.DotStyle)
	// Iterate through each field and add it to the bluge document
	for key, value := range flatDoc {
		if value == nil {
			continue
		}

		if _, ok := mappings.Properties[key]; !ok {
			// Use reflection to find the type of the value.
			// Bluge requires the field type to be specified.
			v := reflect.ValueOf(value)

			// try to find the type of the value and use it to define default mapping
			switch v.Type().String() {
			case "string":
				mappings.Properties[key] = meta.NewProperty("text")
			case "float64":
				mappings.Properties[key] = meta.NewProperty("numeric")
			case "bool":
				mappings.Properties[key] = meta.NewProperty("bool")
			case "time.Time":
				mappings.Properties[key] = meta.NewProperty("time")
			}

			mappingsNeedsUpdate = true
		}

		if !mappings.Properties[key].Index {
			continue // not index, skip
		}

		var field *bluge.TermField
		switch mappings.Properties[key].Type {
		case "text":
			field = bluge.NewTextField(key, value.(string)).SearchTermPositions()
			fieldAnalyzer, _ := zincanalysis.QueryAnalyzerForField(index.CachedAnalyzers, index.CachedMappings, key)
			if fieldAnalyzer != nil {
				field.WithAnalyzer(fieldAnalyzer)
			}
		case "numeric":
			field = bluge.NewNumericField(key, value.(float64))
		case "keyword":
			// compatible verion <= v0.1.4
			if v, ok := value.(bool); ok {
				field = bluge.NewKeywordField(key, strconv.FormatBool(v))
			} else if v, ok := value.(string); ok {
				field = bluge.NewKeywordField(key, v)
			} else {
				return nil, fmt.Errorf("keyword type only support text")
			}
		case "bool": // found using existing index mapping
			value := value.(bool)
			field = bluge.NewKeywordField(key, strconv.FormatBool(value))
		case "time":
			format := time.RFC3339
			if mappings.Properties[key].Format != "" {
				format = mappings.Properties[key].Format
			}
			tim, err := time.Parse(format, value.(string))
			if err != nil {
				return nil, err
			}
			field = bluge.NewDateTimeField(key, tim)
		}

		if mappings.Properties[key].Store {
			field.StoreValue()
		}
		if mappings.Properties[key].Sortable {
			field.Sortable()
		}
		if mappings.Properties[key].Aggregatable {
			field.Aggregatable()
		}
		if mappings.Properties[key].Highlightable {
			field.HighlightMatches()
		}
		bdoc.AddField(field)
	}

	if mappingsNeedsUpdate {
		index.SetMappings(mappings)
		StoreIndex(index)
	}

	docByteVal, _ := json.Marshal(*doc)
	bdoc.AddField(bluge.NewDateTimeField("@timestamp", time.Now()).StoreValue().Sortable().Aggregatable())
	bdoc.AddField(bluge.NewStoredOnlyField("_index", []byte(index.Name)))
	bdoc.AddField(bluge.NewStoredOnlyField("_source", docByteVal))
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil)) // Add _all field that can be used for search

	return bdoc, nil
}

func (index *Index) UseTemplate() error {
	template, err := UseTemplate(index.Name)
	if err != nil {
		return err
	}

	if template == nil {
		return nil
	}

	if template.Template.Settings != nil {
		index.SetSettings(template.Template.Settings)
	}

	if template.Template.Mappings != nil {
		index.SetMappings(template.Template.Mappings)
	}

	return nil
}

func (index *Index) SetSettings(settings *meta.IndexSettings) error {
	if settings == nil {
		return nil
	}

	if settings.NumberOfShards == 0 {
		settings.NumberOfShards = 3
	}
	if settings.NumberOfReplicas == 0 {
		settings.NumberOfReplicas = 1
	}

	index.Settings = settings

	return nil
}

func (index *Index) SetAnalyzers(analyzers map[string]*analysis.Analyzer) error {
	if len(analyzers) == 0 {
		return nil
	}

	index.CachedAnalyzers = analyzers

	return nil
}

func (index *Index) SetMappings(mappings *meta.Mappings) error {
	if mappings == nil || len(mappings.Properties) == 0 {
		return nil
	}

	// custom analyzer just for text field
	for _, prop := range mappings.Properties {
		if prop.Type != "text" {
			prop.Analyzer = ""
			prop.SearchAnalyzer = ""
		}
	}

	mappings.Properties["_id"] = meta.NewProperty("keyword")

	// @timestamp need date_range/date_histogram aggregation, and mappings used for type check in aggregation
	mappings.Properties["@timestamp"] = meta.NewProperty("time")

	// update in the cache
	index.CachedMappings = mappings
	index.Mappings = nil

	return nil
}

// DEPRECATED GetStoredMapping returns the mappings of all the indexes from _index_mapping system index
func (index *Index) GetStoredMapping() (*meta.Mappings, error) {
	log.Error().Bool("deprecated", true).Msg("GetStoredMapping is deprecated, use index.CachedMappings instead")
	for _, indexName := range systemIndexList {
		if index.Name == indexName {
			return nil, nil
		}
	}

	reader, _ := ZINC_SYSTEM_INDEX_LIST["_index_mapping"].Writer.Reader()
	defer reader.Close()

	// search for the index mapping _index_mapping index
	query := bluge.NewTermQuery(index.Name).SetField("_id")
	searchRequest := bluge.NewTopNSearch(1, query) // Should get just 1 result at max
	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Error().Str("index", index.Name).Msg("error executing search: " + err.Error())
		return nil, err
	}

	next, err := dmi.Next()
	if err != nil {
		return nil, err
	}

	mappings := new(meta.Mappings)
	oldMappings := make(map[string]string)
	if next != nil {
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			switch field {
			case "_source":
				if string(value) != "" {
					json.Unmarshal(value, mappings)
				}
			default:
				oldMappings[field] = string(value)
			}
			return true
		})
		if err != nil {
			return nil, err
		}
	}

	// compatible old mappings format
	if len(mappings.Properties) == 0 && len(oldMappings) > 0 {
		mappings.Properties = make(map[string]meta.Property, len(oldMappings))
		for k, v := range oldMappings {
			mappings.Properties[k] = meta.NewProperty(v)
		}
	}

	if len(mappings.Properties) == 0 {
		mappings.Properties = make(map[string]meta.Property)
	}

	return mappings, nil
}
