/* Copyright 2022 Zinc Labs Inc. and Contributors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package core

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/goccy/go-json"

	"github.com/zinclabs/zinc/pkg/meta"
	zincanalysis "github.com/zinclabs/zinc/pkg/uquery/analysis"
	"github.com/zinclabs/zinc/pkg/zutils"
	"github.com/zinclabs/zinc/pkg/zutils/flatten"
)

var errIdNotFound = errors.New("id not found")

// BuildBlugeDocumentFromJSON returns the bluge document for the json document. It also updates the mapping for the fields if not found.
// If no mappings are found, it creates te mapping for all the encountered fields. If mapping for some fields is found but not for others
// then it creates the mapping for the missing fields.
func (index *Index) BuildBlugeDocumentFromJSON(docID string, doc map[string]interface{}) (*bluge.Document, error) {
	// Pick the index mapping from the cache if it already exists
	mappings := index.Mappings
	if mappings == nil {
		mappings = meta.NewMappings()
	}

	mappingsNeedsUpdate := false

	// Create a new bluge document
	bdoc := bluge.NewDocument(docID)
	flatDoc, _ := flatten.Flatten(doc, "")
	// Iterate through each field and add it to the bluge document
	for key, value := range flatDoc {
		if value == nil || key == meta.TimeFieldName {
			continue
		}

		if _, ok := mappings.GetProperty(key); !ok {
			// try to find the type of the value and use it to define default mapping
			switch value.(type) {
			case string:
				mappings.SetProperty(key, meta.NewProperty("text"))
			case int, int64, float64:
				mappings.SetProperty(key, meta.NewProperty("numeric"))
			case bool:
				mappings.SetProperty(key, meta.NewProperty("bool"))
			case []interface{}:
				if v, ok := value.([]interface{}); ok {
					for _, vv := range v {
						switch vv.(type) {
						case string:
							mappings.SetProperty(key, meta.NewProperty("text"))
						case float64:
							mappings.SetProperty(key, meta.NewProperty("numeric"))
						case bool:
							mappings.SetProperty(key, meta.NewProperty("bool"))
						}
						break
					}
				}
			}

			mappingsNeedsUpdate = true
		}

		if prop, ok := mappings.GetProperty(key); ok && !prop.Index {
			continue // not index, skip
		}

		switch v := value.(type) {
		case []interface{}:
			for _, v := range v {
				if err := index.buildField(mappings, bdoc, key, v); err != nil {
					return nil, err
				}
			}
		default:
			if err := index.buildField(mappings, bdoc, key, v); err != nil {
				return nil, err
			}
		}
	}

	if mappingsNeedsUpdate {
		_ = index.SetMappings(mappings)
		_ = StoreIndex(index)
	}

	// set timestamp
	timestamp := time.Now()
	if value, ok := flatDoc[meta.TimeFieldName]; ok {
		delete(doc, meta.TimeFieldName)
		prop, _ := mappings.GetProperty(meta.TimeFieldName)
		v, err := zutils.ParseTime(value, prop.Format, prop.TimeZone)
		if err != nil {
			return nil, fmt.Errorf("field [%s] value [%v] parse err: %s", meta.TimeFieldName, value, err.Error())
		}
		timestamp = v
	}
	bdoc.AddField(bluge.NewDateTimeField(meta.TimeFieldName, timestamp).StoreValue().Sortable().Aggregatable())

	docByteVal, _ := json.Marshal(doc)
	bdoc.AddField(bluge.NewStoredOnlyField("_index", []byte(index.Name)))
	bdoc.AddField(bluge.NewStoredOnlyField("_source", docByteVal))
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", []string{"_id", "_index", "_source", meta.TimeFieldName}))

	// Add time for index
	bdoc.SetTimestamp(timestamp.UnixNano())
	// Upate metadata
	index.SetTimestamp(timestamp.UnixNano())

	return bdoc, nil
}

func (index *Index) buildField(mappings *meta.Mappings, bdoc *bluge.Document, key string, value interface{}) error {
	var field *bluge.TermField
	prop, _ := mappings.GetProperty(key)
	switch prop.Type {
	case "text":
		v, err := zutils.ToString(value)
		if err != nil {
			return fmt.Errorf("field [%s] was set type to [text] but the value [%v] can't convert to string", key, value)
		}
		field = bluge.NewTextField(key, v).SearchTermPositions()
		fieldAnalyzer, _ := zincanalysis.QueryAnalyzerForField(index.Analyzers, index.Mappings, key)
		if fieldAnalyzer != nil {
			field.WithAnalyzer(fieldAnalyzer)
		}
	case "numeric":
		v, err := zutils.ToFloat64(value)
		if err != nil {
			return fmt.Errorf("field [%s] was set type to [numeric] but the value [%v] can't convert to int", key, value)
		}
		field = bluge.NewNumericField(key, v)
	case "keyword":
		v, err := zutils.ToString(value)
		if err != nil {
			return fmt.Errorf("field [%s] was set type to [keyword] but the value [%v] can't convert to string", key, value)
		}
		field = bluge.NewKeywordField(key, v)
	case "bool":
		v, err := zutils.ToBool(value)
		if err != nil {
			return fmt.Errorf("field [%s] was set type to [bool] but the value [%v] can't convert to boolean", key, value)
		}
		field = bluge.NewKeywordField(key, strconv.FormatBool(v))
	case "date", "time":
		v, err := zutils.ParseTime(value, prop.Format, prop.TimeZone)
		if err != nil {
			return fmt.Errorf("field [%s] value [%v] parse err: %s", key, value, err.Error())
		}
		field = bluge.NewDateTimeField(key, v)
	}
	if prop.Store || prop.Highlightable {
		field.StoreValue()
	}
	if prop.Highlightable {
		field.HighlightMatches()
	}
	if prop.Sortable {
		field.Sortable()
	}
	if prop.Aggregatable {
		field.Aggregatable()
	}
	bdoc.AddField(field)

	return nil
}

func (index *Index) CreateDocumentAsync(docID string, doc map[string]interface{}) error {
	return index.wal.CreateDocument(docID, doc)
}

func (index *Index) UpdateDocumentAsync(docID string, doc map[string]interface{}) error {
	return index.wal.UpdateDocument(docID, doc)
}

func (index *Index) DeleteDocumentAsync(docID string) error {
	return index.wal.DeleteDocument(docID)
}

// CreateDocument inserts or updates a document in the zinc index
func (index *Index) CreateDocument(docID string, doc map[string]interface{}, update bool) error {
	bdoc, err := index.BuildBlugeDocumentFromJSON(docID, doc)
	if err != nil {
		return err
	}

	// Finally update the document on disk
	writer, err := index.GetWriter()
	if err != nil {
		return err
	}
	if update {
		err = writer.Update(bdoc.ID(), bdoc)
	} else {
		err = writer.Insert(bdoc)
	}
	return err
}

// UpdateDocument updates a document in the zinc index
func (index *Index) UpdateDocument(docID string, doc map[string]interface{}) error {
	bdoc, err := index.BuildBlugeDocumentFromJSON(docID, doc)
	if err != nil {
		return err
	}
	writer, err := index.FindID(docID)
	if err != nil {
		return err
	}
	return writer.Update(bdoc.ID(), bdoc)
}

func (index *Index) DeleteDocument(docID string) error {
	bdoc := bluge.NewDocument(docID)
	writers, err := index.GetWriters()
	if err != nil {
		return err
	}
	for _, w := range writers {
		err = w.Delete(bdoc.ID())
		if err != nil {
			return err
		}
	}
	return nil
}

func (index *Index) FindID(id string) (*bluge.Writer, error) {
	query := bluge.NewBooleanQuery()
	query.AddMust(bluge.NewTermQuery(id).SetField("_id"))
	request := bluge.NewTopNSearch(1, query).WithStandardAggregations()
	ctx := context.Background()

	// check id store by which shard
	writers, err := index.GetWriters()
	if err != nil {
		return nil, err
	}

	for _, w := range writers {
		r, err := w.Reader()
		if err != nil {
			return nil, err
		}
		defer r.Close()
		dmi, err := r.Search(ctx, request)
		if err != nil {
			return nil, err
		}
		if dmi.Aggregations().Count() > 0 {
			return w, nil
		}
	}
	return nil, errIdNotFound
}
