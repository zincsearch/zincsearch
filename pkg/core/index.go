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
	"github.com/prabhatsharma/zinc/pkg/zutils"
	"github.com/rs/zerolog/log"
)

// BuildBlugeDocumentFromJSON returns the bluge document for the json document. It also updates the mapping for the fields if not found.
// If no mappings are found, it creates te mapping for all the encountered fields. If mapping for some fields is found but not for others
// then it creates the mapping for the missing fields.
func (rindex *Index) BuildBlugeDocumentFromJSON(docID string, doc *map[string]interface{}) (*bluge.Document, error) {

	// Pick the index mapping from the cache if it already exists
	indexMapping := rindex.CachedMapping

	if indexMapping == nil {
		indexMapping = make(map[string]string)
	}

	flatDoc, _ := flatten.Flatten(*doc, "", flatten.DotStyle)

	// Create a new bluge document
	bdoc := bluge.NewDocument(docID)

	indexMappingNeedsUpdate := false

	// Iterate through each field and add it to the bluge document
	for key, value := range flatDoc {
		if _, ok := indexMapping[key]; !ok {
			// Assign auto inferred type for the new key

			// Use reflection to find the type of the value.
			// Bluge requires the field type to be specified.

			if value != nil { // value could be just {} in the json data or e.g. "rules": null or "creationTimestamp": null, etc.
				v := reflect.ValueOf(value)

				// try to find the type of the value and use it to define default mapping
				switch v.Type().String() {
				case "string":
					indexMapping[key] = "text"
				case "float64":
					indexMapping[key] = "numeric"
				case "bool":
					indexMapping[key] = "bool"
				case "time.Time":
					indexMapping[key] = "time"
				}

				indexMappingNeedsUpdate = true
			}
		}

		if value != nil {
			switch indexMapping[key] {
			case "text":
				stringField := bluge.NewTextField(key, value.(string)).SearchTermPositions().Aggregatable()
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

	if indexMappingNeedsUpdate {
		rindex.SetMapping(indexMapping)
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
func (index *Index) SetMapping(iMap map[string]string) error {
	bdoc := bluge.NewDocument(index.Name)

	// @timestamp need date_range/date_histogram aggregation, and mappings used for type check in aggregation
	iMap["@timestamp"] = "time"

	for k, v := range iMap {
		bdoc.AddField(bluge.NewTextField(k, v).StoreValue())
	}

	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil))

	// update on the disk
	systemIndex := ZINC_SYSTEM_INDEX_LIST["_index_mapping"].Writer
	err := systemIndex.Update(bdoc.ID(), bdoc)
	if err != nil {
		log.Printf("error updating document: %v", err)
		return err
	}

	// update in the cache
	index.CachedMapping = iMap

	return nil
}

// GetStoredMapping returns the mappings of all the indexes from _index_mapping system index
func (index *Index) GetStoredMapping() (map[string]string, error) {
	dataPath := zutils.GetEnv("ZINC_DATA_PATH", "./data")
	systemPath := dataPath + "/_index_mapping"

	config := bluge.DefaultConfig(systemPath)
	reader, err := bluge.OpenReader(config)
	if err != nil {
		return nil, nil //probably no system index available
		// log.Fatalf("GetIndexMapping: unable to open reader: %v", err)
	}
	defer reader.Close()

	// search for the index mapping _index_mapping index
	query := bluge.NewTermQuery(index.Name).SetField("_id")
	searchRequest := bluge.NewTopNSearch(1, query) // Should get just 1 result at max

	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Log().Msg("error executing search: " + err.Error())
	}

	next, err := dmi.Next()
	if err != nil {
		return nil, err
	}

	if next != nil {
		result := make(map[string]string)
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			result[field] = string(value)
			return true
		})
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, nil
}
