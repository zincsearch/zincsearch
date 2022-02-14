package core

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/jeremywohl/flatten"
	"github.com/rs/zerolog/log"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

// BuildBlugeDocumentFromJSON returns the bluge document for the json document. It also updates the mapping for the fields if not found.
// If no mappings are found, it creates te mapping for all the encountered fields. If mapping for some fields is found but not for others
// then it creates the mapping for the missing fields.
func (index *Index) BuildBlugeDocumentFromJSON(docID string, doc *map[string]interface{}) (*bluge.Document, error) {
	// Pick the index mapping from the cache if it already exists
	mappings := index.CachedMappings
	if mappings == nil {
		mappings = new(meta.Mappings)
		mappings.Properties = make(map[string]meta.Property)
	}

	mappingsNeedsUpdate := false

	// Create a new bluge document
	bdoc := bluge.NewDocument(docID)
	flatDoc, _ := flatten.Flatten(*doc, "", flatten.DotStyle)
	// Iterate through each field and add it to the bluge document
	for key, value := range flatDoc {
		if _, ok := mappings.Properties[key]; !ok {
			// Assign auto inferred type for the new key

			// Use reflection to find the type of the value.
			// Bluge requires the field type to be specified.

			if value != nil { // value could be just {} in the json data or e.g. "rules": null or "creationTimestamp": null, etc.
				v := reflect.ValueOf(value)

				// try to find the type of the value and use it to define default mapping
				switch v.Type().String() {
				case "string":
					mappings.Properties[key] = meta.Property{Type: "text"}
				case "float64":
					mappings.Properties[key] = meta.Property{Type: "numeric"}
				case "bool":
					mappings.Properties[key] = meta.Property{Type: "bool"}
				case "time.Time":
					mappings.Properties[key] = meta.Property{Type: "time"}
				}

				mappingsNeedsUpdate = true
			}
		}

		if value != nil {
			switch mappings.Properties[key].Type {
			case "text":
				stringField := bluge.NewTextField(key, value.(string)).SearchTermPositions().Aggregatable().StoreValue().HighlightMatches()
				bdoc.AddField(stringField)
			case "numeric":
				numericField := bluge.NewNumericField(key, value.(float64)).Aggregatable()
				bdoc.AddField(numericField)
			case "keyword":
				// compatible verion <= v0.1.4
				var keywordField *bluge.TermField
				if v, ok := value.(bool); ok {
					keywordField = bluge.NewKeywordField(key, strconv.FormatBool(v)).Aggregatable()
				} else if v, ok := value.(string); ok {
					keywordField = bluge.NewKeywordField(key, v).Aggregatable()
				} else {
					return nil, fmt.Errorf("keyword type only support text")
				}
				bdoc.AddField(keywordField)
			case "bool": // found using existing index mapping
				value := value.(bool)
				keywordField := bluge.NewKeywordField(key, strconv.FormatBool(value)).Aggregatable()
				bdoc.AddField(keywordField)
			case "time":
				tim, err := time.Parse(time.RFC3339, value.(string))
				if err != nil {
					return nil, err
				}
				timeField := bluge.NewDateTimeField(key, tim).Aggregatable()
				bdoc.AddField(timeField)
			}
		}
	}

	if mappingsNeedsUpdate {
		index.SetMappings(mappings)
	}

	docByteVal, _ := json.Marshal(*doc)
	bdoc.AddField(bluge.NewDateTimeField("@timestamp", time.Now()).StoreValue().Aggregatable())
	bdoc.AddField(bluge.NewStoredOnlyField("_source", docByteVal))
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil)) // Add _all field that can be used for search

	return bdoc, nil
}

// SetMapping Saves the mapping of the index to _index_mapping index
// index: Name of the index ffor which the mapping needs to be saved
// iMap: a map of the fileds that specify name and type of the field. e.g. movietitle: string
func (index *Index) SetMappings(mappings *meta.Mappings) error {
	// @timestamp need date_range/date_histogram aggregation, and mappings used for type check in aggregation
	mappings.Properties["@timestamp"] = meta.Property{Type: "time"}

	bdoc := bluge.NewDocument(index.Name)
	for k, prop := range mappings.Properties {
		bdoc.AddField(bluge.NewTextField(k, prop.Type).StoreValue())
	}

	docByteVal, _ := json.Marshal(mappings)
	bdoc.AddField(bluge.NewStoredOnlyField("_source", docByteVal))
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil))

	// update on the disk
	systemIndex := ZINC_SYSTEM_INDEX_LIST["_index_mapping"].Writer
	err := systemIndex.Update(bdoc.ID(), bdoc)
	if err != nil {
		log.Printf("error updating document: %v", err)
		return err
	}

	// update in the cache
	index.CachedMappings = mappings

	return nil
}

// GetStoredMappings returns the mappings of all the indexes from _index_mapping system index
func (index *Index) GetStoredMappings() (*meta.Mappings, error) {
	DATA_PATH := zutils.GetEnv("ZINC_DATA_PATH", "./data")
	systemPath := DATA_PATH + "/_index_mapping"

	config := bluge.DefaultConfig(systemPath)
	reader, err := bluge.OpenReader(config)
	if err != nil {
		log.Error().Str("index", index.Name).Msgf("GetIndexMapping: unable to open reader: %v", err)
		return nil, nil
	}
	defer reader.Close()

	// search for the index mapping _index_mapping index
	query := bluge.NewTermQuery(index.Name).SetField("_id")
	searchRequest := bluge.NewTopNSearch(1, query) // Should get just 1 result at max
	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Error().Str("index", index.Name).Msg("error executing search: " + err.Error())
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
	if len(oldMappings) > 0 && len(mappings.Properties) == 0 {
		mappings.Properties = make(map[string]meta.Property, len(oldMappings))
		for k, v := range oldMappings {
			mappings.Properties[k] = meta.Property{Type: v}
		}
	}

	return mappings, nil
}
