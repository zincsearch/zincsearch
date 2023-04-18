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
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/metadata"
)

// ListTemplates returns all templates
func ListTemplates(pattern string) ([]*meta.Template, error) {
	templates, err := metadata.Template.List(0, 0)
	if err != nil {
		return nil, err
	}
	if templates == nil {
		templates = make([]*meta.Template, 0)
	}
	if pattern != "" {
		oldTpls := templates[:]
		templates = templates[:0]
		for _, tpl := range oldTpls {
			for i := range tpl.IndexTemplate.IndexPatterns {
				if tpl.IndexTemplate.IndexPatterns[i] == pattern {
					templates = append(templates, tpl)
					break
				}
			}
		}
	}
	return templates, nil
}

// NewTemplate create a template and store in local
func NewTemplate(name string, template *meta.IndexTemplate) error {
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

	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	tpl := meta.Template{
		Name:          name,
		IndexTemplate: template,
	}
	err := metadata.Template.Set(name, tpl)
	if err != nil {
		return fmt.Errorf("template: error updating document: %s", err.Error())
	}

	return nil
}

// LoadTemplate load a specific template from local
func LoadTemplate(name string) (*meta.IndexTemplate, bool, error) {
	if name == "" {
		return nil, false, nil
	}

	tpl, err := metadata.Template.Get(name)
	if err != nil {
		if err == errors.ErrKeyNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return tpl.IndexTemplate, true, nil
}

// DeleteTemplate delete a template from local
func DeleteTemplate(name string) error {
	return metadata.Template.Delete(name)
}

// UseTemplate use a specific template for new index
func UseTemplate(indexName string) (*meta.IndexTemplate, error) {
	templates, err := ListTemplates("")
	if err != nil {
		return nil, err
	}
	if len(templates) == 0 {
		return nil, nil
	}

	// sort by priority
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].IndexTemplate.Priority > templates[j].IndexTemplate.Priority
	})

	// filter by first character
	var filteredTemplates []*meta.Template
	for i := range templates {
		for j := range templates[i].IndexTemplate.IndexPatterns {
			if strings.HasPrefix(indexName, templates[i].IndexTemplate.IndexPatterns[j][:1]) {
				filteredTemplates = append(filteredTemplates, templates[i])
				break
			}
		}
	}

	for _, tpl := range filteredTemplates {
		for _, pattern := range tpl.IndexTemplate.IndexPatterns {
			pattern := strings.TrimRight(strings.ReplaceAll(pattern, "*", ".*"), "$") + "$"
			re := regexp.MustCompile(pattern)
			if re.MatchString(indexName) {
				return tpl.IndexTemplate, nil
			}
		}
	}

	return nil, nil
}
