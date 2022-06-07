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
	"fmt"
	"strconv"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/goccy/go-json"

	"github.com/zinclabs/zinc/pkg/meta"
	zincanalysis "github.com/zinclabs/zinc/pkg/uquery/analysis"
	"github.com/zinclabs/zinc/pkg/zutils"
	"github.com/zinclabs/zinc/pkg/zutils/flatten"
)

type Index struct {
	meta.Index
	Analyzers map[string]*analysis.Analyzer `json:"-"`
	Writer    *bluge.Writer                 `json:"-"`
}

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
		if value == nil || key == "@timestamp" {
			continue
		}

		if _, ok := mappings.GetProperty(key); !ok {
			// try to find the type of the value and use it to define default mapping
			switch value.(type) {
			case string:
				mappings.SetProperty(key, meta.NewProperty("text"))
			case float64:
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

	timestamp := time.Now()
	if v, ok := flatDoc["@timestamp"]; ok {
		switch v := v.(type) {
		case string:
			if t, err := time.Parse(time.RFC3339, v); err == nil && !t.IsZero() {
				timestamp = t
				delete(doc, "@timestamp")
			}
		case float64:
			if t := zutils.Unix(int64(v)); !t.IsZero() {
				timestamp = t
				delete(doc, "@timestamp")
			}
		default:
			// noop
		}
	}
	docByteVal, _ := json.Marshal(doc)
	bdoc.AddField(bluge.NewDateTimeField("@timestamp", timestamp).StoreValue().Sortable().Aggregatable())
	bdoc.AddField(bluge.NewStoredOnlyField("_index", []byte(index.Name)))
	bdoc.AddField(bluge.NewStoredOnlyField("_source", docByteVal))
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", []string{"_id", "_index", "_source", "@timestamp"}))

	// test for add time index
	bdoc.SetTimestamp(timestamp.UnixNano())

	return bdoc, nil
}

func (index *Index) buildField(mappings *meta.Mappings, bdoc *bluge.Document, key string, value interface{}) error {
	var field *bluge.TermField
	prop, _ := mappings.GetProperty(key)
	switch prop.Type {
	case "text":
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("field [%s] was set type to [text] but got a %T value", key, value)
		}
		field = bluge.NewTextField(key, v).SearchTermPositions()
		fieldAnalyzer, _ := zincanalysis.QueryAnalyzerForField(index.Analyzers, index.Mappings, key)
		if fieldAnalyzer != nil {
			field.WithAnalyzer(fieldAnalyzer)
		}
	case "numeric":
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("field [%s] was set type to [numeric] but got a %T value", key, value)
		}
		field = bluge.NewNumericField(key, v)
	case "keyword":
		switch v := value.(type) {
		case string:
			field = bluge.NewKeywordField(key, v)
		case float64:
			field = bluge.NewKeywordField(key, strconv.FormatFloat(v, 'f', -1, 64))
		case int:
			field = bluge.NewKeywordField(key, strconv.FormatInt(int64(v), 10))
		case bool:
			field = bluge.NewKeywordField(key, strconv.FormatBool(v))
		default:
			field = bluge.NewKeywordField(key, fmt.Sprintf("%v", v))
		}
	case "bool":
		value := value.(bool)
		field = bluge.NewKeywordField(key, strconv.FormatBool(value))
	case "date", "time":
		switch v := value.(type) {
		case string:
			format := time.RFC3339
			if prop.Format != "" {
				format = prop.Format
			}
			var tim time.Time
			var err error
			if format == "epoch_millis" {
				tim = time.UnixMilli(int64(value.(float64)))
			} else {
				tim, err = time.Parse(format, value.(string))
			}
			if err != nil {
				return err
			}
			field = bluge.NewDateTimeField(key, tim)
		case float64:
			if t := zutils.Unix(int64(v)); !t.IsZero() {
				field = bluge.NewDateTimeField(key, t)
			}
		}
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

func (index *Index) UseTemplate() error {
	template, err := UseTemplate(index.Name)
	if err != nil {
		return err
	}

	if template == nil {
		return nil
	}

	if template.Template.Settings != nil {
		_ = index.SetSettings(template.Template.Settings)
	}

	if template.Template.Mappings != nil {
		_ = index.SetMappings(template.Template.Mappings)
	}

	return nil
}

func (index *Index) SetSettings(settings *meta.IndexSettings) error {
	if settings == nil {
		return nil
	}

	index.Settings = settings

	return nil
}

func (index *Index) SetAnalyzers(analyzers map[string]*analysis.Analyzer) error {
	if len(analyzers) == 0 {
		return nil
	}

	index.Analyzers = analyzers

	return nil
}

func (index *Index) SetMappings(mappings *meta.Mappings) error {
	if mappings == nil || mappings.Len() == 0 {
		return nil
	}

	// custom analyzer just for text field
	for field, prop := range mappings.ListProperty() {
		if prop.Type != "text" {
			prop.Analyzer = ""
			prop.SearchAnalyzer = ""
			mappings.SetProperty(field, prop)
		}
	}

	mappings.SetProperty("_id", meta.NewProperty("keyword"))

	// @timestamp need date_range/date_histogram aggregation, and mappings used for type check in aggregation
	mappings.SetProperty("@timestamp", meta.NewProperty("date"))

	// update in the cache
	index.Mappings = mappings

	return nil
}

func (index *Index) UpdateMetadata() {
	w := index.Writer
	if w == nil {
		return
	}
	status := w.Status()
	index.StorageSize = status.CurOnDiskBytes

	if r, err := w.Reader(); err == nil {
		if n, err := r.Count(); err == nil {
			index.DocNum = n
		}
	}
}

func (index *Index) Close() error {
	// update metadata before close
	index.UpdateMetadata()
	return index.Writer.Close()
}
