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

package meta

import (
	"bytes"
	"encoding/json"
	"sync"
)

type Mappings struct {
	Properties map[string]Property `json:"properties,omitempty"`
	lock       sync.RWMutex
}

type Property struct {
	Type           string `json:"type"` // text, keyword, date, numeric, boolean, geo_point
	Analyzer       string `json:"analyzer,omitempty"`
	SearchAnalyzer string `json:"search_analyzer,omitempty"`
	Format         string `json:"format,omitempty"` // date format yyyy-MM-dd HH:mm:ss || yyyy-MM-dd || epoch_millis
	Index          bool   `json:"index"`
	Store          bool   `json:"store"`
	Sortable       bool   `json:"sortable"`
	Aggregatable   bool   `json:"aggregatable"`
	Highlightable  bool   `json:"highlightable"`
}

func NewMappings() *Mappings {
	return &Mappings{
		Properties: make(map[string]Property),
	}
}

func NewProperty(typ string) Property {
	p := Property{
		Type:           typ,
		Analyzer:       "",
		SearchAnalyzer: "",
		Format:         "",
		Index:          true,
		Store:          false,
		Sortable:       true,
		Aggregatable:   true,
		Highlightable:  false,
	}
	if typ == "text" {
		p.Sortable = false
		p.Aggregatable = false
		p.Store = true
		p.Highlightable = true
	}
	return p
}

func (t *Mappings) Len() int {
	t.lock.RLock()
	n := len(t.Properties)
	t.lock.RUnlock()
	return n
}

func (t *Mappings) SetProperty(field string, prop Property) {
	t.lock.Lock()
	t.Properties[field] = prop
	t.lock.Unlock()
}

func (t *Mappings) GetProperty(field string) (Property, bool) {
	t.lock.RLock()
	prop, ok := t.Properties[field]
	t.lock.RUnlock()
	return prop, ok
}

func (t *Mappings) ListProperty() map[string]Property {
	m := make(map[string]Property)
	t.lock.RLock()
	for k, v := range t.Properties {
		m[k] = v
	}
	t.lock.RUnlock()
	return m
}

func (t *Mappings) MarshalJSON() ([]byte, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	b := bytes.NewBuffer(nil)
	b.WriteString(`{"properties":`)
	p, err := json.Marshal(t.Properties)
	if err != nil {
		return nil, err
	}
	b.Write(p)
	b.WriteByte('}')
	return b.Bytes(), nil
}
