package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/rs/zerolog/log"

	"github.com/prabhatsharma/zinc/pkg/dsl/meta"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
)

// ListTemplates returns all templates
func ListTemplates() (*v1.SearchResponse, error) {
	query := bluge.NewMatchAllQuery()
	searchRequest := bluge.NewTopNSearch(1000, query).WithStandardAggregations()
	reader, _ := ZINC_SYSTEM_INDEX_LIST["_index_template"].Writer.Reader()
	defer reader.Close()

	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		return nil, fmt.Errorf("core.ListTemplates: error executing search: %v", err)
	}

	var Hits []v1.Hit
	next, err := dmi.Next()
	for err == nil && next != nil {
		var id string
		var timestamp time.Time
		tpl := new(IndexTemplate)
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			switch field {
			case "_id":
				id = string(value)
			case "name":
				tpl.Name = string(value)
			case "priority":
				priority, _ := bluge.DecodeNumericFloat64(value)
				tpl.Priority = int(priority)
			case "index_prefix":
				tpl.IndexPrefix = string(value)
			case "@timestamp":
				timestamp, _ = bluge.DecodeDateTime(value)
			default:
				if strings.HasPrefix(field, "index_pattern_") {
					tpl.IndexPatterns = append(tpl.IndexPatterns, string(value))
				}
			}

			return true
		})
		if err != nil {
			log.Printf("core.ListTemplates: error accessing stored fields: %v", err)
		}

		hit := v1.Hit{
			Index:     tpl.Name,
			Type:      tpl.Name,
			ID:        id,
			Score:     next.Score,
			Timestamp: timestamp,
			Source:    tpl,
		}
		Hits = append(Hits, hit)

		next, err = dmi.Next()
	}

	resp := &v1.SearchResponse{
		Took: int(dmi.Aggregations().Duration().Milliseconds()),
		Hits: v1.Hits{
			Total: v1.Total{
				Value: int(dmi.Aggregations().Count()),
			},
			MaxScore: dmi.Aggregations().Metric("max_score"),
			Hits:     Hits,
		},
	}

	return resp, nil
}

// NewTemplate create a template and store in local
func NewTemplate(name string, template *meta.Template) error {
	update := false
	_, tplExists, _ := LoadTemplate(name)
	if tplExists {
		update = true
	}

	bdoc := bluge.NewDocument(name)
	bdoc.AddField(bluge.NewKeywordField("name", name).StoreValue().Sortable())
	bdoc.AddField(bluge.NewNumericField("priority", float64(template.Priority)).StoreValue().Sortable().Aggregatable())
	for i := 0; i < len(template.IndexPatterns); i++ {
		name := fmt.Sprintf("index_pattern_%d", i)
		bdoc.AddField(bluge.NewKeywordField(name, template.IndexPatterns[i]).StoreValue())
		bdoc.AddField(bluge.NewKeywordField("index_prefix", string(template.IndexPatterns[i][0:1])).StoreValue())
	}

	docByteVal, _ := json.Marshal(*template)
	bdoc.AddField(bluge.NewDateTimeField("@timestamp", time.Now()).StoreValue().Sortable().Aggregatable())
	bdoc.AddField(bluge.NewStoredOnlyField("_source", docByteVal))
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil)) // Add _all field that can be used for search

	var err error
	index := ZINC_SYSTEM_INDEX_LIST["_index_template"].Writer
	if update {
		err = index.Update(bdoc.ID(), bdoc)
	} else {
		err = index.Insert(bdoc)
	}
	if err != nil {
		return fmt.Errorf("template: error updating document: %v", err)
	}

	return nil
}

// LoadTemplate load a specific template from local
func LoadTemplate(name string) (*meta.Template, bool, error) {
	if name == "" {
		return nil, false, nil
	}

	query := bluge.NewTermQuery(name).SetField("_id")
	searchRequest := bluge.NewTopNSearch(1, query)
	reader, _ := ZINC_SYSTEM_INDEX_LIST["_index_template"].Writer.Reader()
	defer reader.Close()

	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		return nil, false, fmt.Errorf("template: error executing search: %v", err)
	}

	tpl := new(meta.Template)
	next, err := dmi.Next()
	if err != nil {
		return nil, false, fmt.Errorf("template: error accessing stored fields: %v", err)
	}
	if next == nil {
		return nil, false, fmt.Errorf("template: %s not found", name)
	}
	err = next.VisitStoredFields(func(field string, value []byte) bool {
		switch field {
		case "_source":
			json.Unmarshal(value, tpl)
			return true
		default:
		}
		return true
	})
	if err != nil {
		return nil, false, fmt.Errorf("template: error accessing stored fields: %v", err)
	}

	return tpl, true, nil
}

// DeleteTemplate delete a template from local
func DeleteTemplate(name string) error {
	bdoc := bluge.NewDocument(name)
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil))
	err := ZINC_SYSTEM_INDEX_LIST["_index_template"].Writer.Delete(bdoc.ID())
	if err != nil {
		return fmt.Errorf("template: error deleting template: %v", err)
	}

	return nil
}

// UseTemplate use a specific template for new index
func UseTemplate(indexName string) (*meta.Template, error) {
	query := bluge.NewTermQuery(string(indexName[0:1])).SetField("index_prefix")
	searchRequest := bluge.NewTopNSearch(1000, query).SortBy([]string{"-priority"}).WithStandardAggregations()
	reader, _ := ZINC_SYSTEM_INDEX_LIST["_index_template"].Writer.Reader()
	defer reader.Close()

	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		return nil, fmt.Errorf("core.UseTemplate: error executing search: %v", err)
	}

	templates := make([]*meta.Template, 0)
	next, err := dmi.Next()
	for err == nil && next != nil {
		tpl := new(meta.Template)
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			if field == "_source" {
				json.Unmarshal(value, tpl)
			}
			return true
		})
		if err != nil {
			log.Printf("core.UseTemplate: error accessing stored fields: %v", err)
		}

		templates = append(templates, tpl)
		next, err = dmi.Next()
	}

	if err != nil {
		return nil, err
	}

	if len(templates) == 0 {
		return nil, nil
	}

	for _, tpl := range templates {
		for _, pattern := range tpl.IndexPatterns {
			pattern = strings.TrimRight(pattern, "*")
			if strings.HasPrefix(indexName, pattern) {
				return tpl, nil
			}
		}
	}

	return nil, nil
}
