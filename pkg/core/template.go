package core

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"

	meta "github.com/zinclabs/zinc/pkg/meta/v2"
)

// ListTemplates returns all templates
func ListTemplates(pattern string) ([]IndexTemplate, error) {
	var query bluge.Query
	if pattern != "" {
		query = bluge.NewBooleanQuery().AddMust(bluge.NewTermQuery(pattern).SetField("index_pattern"))
	} else {
		query = bluge.NewMatchAllQuery()
	}
	searchRequest := bluge.NewTopNSearch(1000, query).SortBy([]string{"name"})
	reader, _ := ZINC_SYSTEM_INDEX_LIST["_index_template"].Writer.Reader()
	defer reader.Close()

	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		return nil, fmt.Errorf("core.ListTemplates: error executing search: %s", err.Error())
	}

	templates := make([]IndexTemplate, 0)
	next, err := dmi.Next()
	for err == nil && next != nil {
		var name string
		var timestamp time.Time
		tpl := new(meta.Template)
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			switch field {
			case "name":
				name = string(value)
			case "@timestamp":
				timestamp, _ = bluge.DecodeDateTime(value)
			case "_source":
				json.Unmarshal(value, tpl)
			default:
			}
			return true
		})
		if err != nil {
			log.Printf("core.ListTemplates: error accessing stored fields: %s", err.Error())
		}

		templates = append(templates, IndexTemplate{
			Name:          name,
			Timestamp:     timestamp,
			IndexTemplate: tpl,
		})

		next, err = dmi.Next()
	}

	return templates, nil
}

// NewTemplate create a template and store in local
func NewTemplate(name string, template *meta.Template) error {
	if name == "" || template == nil {
		return nil
	}

	// check pattern is exists
	for _, pattern := range template.IndexPatterns {
		results, _ := ListTemplates(pattern)
		for _, result := range results {
			if result.Name == name {
				continue
			}
			if result.IndexTemplate.Priority == template.Priority {
				return fmt.Errorf("index template [%s] has index patterns %s "+
					"matching patterns from existing templates [%s] with patterns (%s => %s) "+
					"that have the same priority [%d], multiple index templates may not match during index creation, "+
					"please use a different priority",
					name, template.IndexPatterns,
					result.Name,
					result.Name, result.IndexTemplate.IndexPatterns,
					template.Priority,
				)
			}
		}
	}

	bdoc := bluge.NewDocument(name)
	bdoc.AddField(bluge.NewKeywordField("name", name).StoreValue().Sortable())
	bdoc.AddField(bluge.NewNumericField("priority", float64(template.Priority)).StoreValue().Sortable().Aggregatable())
	for i := 0; i < len(template.IndexPatterns); i++ {
		bdoc.AddField(bluge.NewKeywordField("index_pattern", template.IndexPatterns[i]).StoreValue())
		bdoc.AddField(bluge.NewKeywordField("index_prefix", string(template.IndexPatterns[i][0:1])).StoreValue())
	}

	docByteVal, _ := json.Marshal(*template)
	bdoc.AddField(bluge.NewDateTimeField("@timestamp", time.Now()).StoreValue().Sortable().Aggregatable())
	bdoc.AddField(bluge.NewStoredOnlyField("_source", docByteVal))
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil)) // Add _all field that can be used for search

	index := ZINC_SYSTEM_INDEX_LIST["_index_template"].Writer
	err := index.Update(bdoc.ID(), bdoc)
	if err != nil {
		return fmt.Errorf("template: error updating document: %s", err.Error())
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
		return nil, false, fmt.Errorf("template: error executing search: %s", err.Error())
	}

	tpl := new(meta.Template)
	next, err := dmi.Next()
	if err != nil {
		return nil, false, fmt.Errorf("template: error accessing stored fields: %s", err.Error())
	}
	if next == nil {
		return nil, false, nil
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
		return nil, false, fmt.Errorf("template: error accessing stored fields: %s", err.Error())
	}

	return tpl, true, nil
}

// DeleteTemplate delete a template from local
func DeleteTemplate(name string) error {
	bdoc := bluge.NewDocument(name)
	bdoc.AddField(bluge.NewCompositeFieldExcluding("_all", nil))
	err := ZINC_SYSTEM_INDEX_LIST["_index_template"].Writer.Delete(bdoc.ID())
	if err != nil {
		return fmt.Errorf("template: error deleting template: %s", err.Error())
	}

	return nil
}

// UseTemplate use a specific template for new index
func UseTemplate(indexName string) (*meta.Template, error) {
	query := bluge.NewTermQuery(string(indexName[0:1])).SetField("index_prefix")
	searchRequest := bluge.NewTopNSearch(1000, query).SortBy([]string{"-priority"})
	reader, _ := ZINC_SYSTEM_INDEX_LIST["_index_template"].Writer.Reader()
	defer reader.Close()

	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		return nil, fmt.Errorf("core.UseTemplate: error executing search: %s", err.Error())
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
			log.Printf("core.UseTemplate: error accessing stored fields: %s", err.Error())
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
			pattern := strings.TrimRight(strings.ReplaceAll(pattern, "*", ".*"), "$") + "$"
			re := regexp.MustCompile(pattern)
			if re.MatchString(indexName) {
				return tpl, nil
			}
		}
	}

	return nil, nil
}
