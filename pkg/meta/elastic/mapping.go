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

package elastic

import (
	"bytes"
	"sync"

	"github.com/zincsearch/zincsearch/pkg/zutils/json"
)

// Mappings holds the index mappings.
type Mappings struct {
	lock sync.RWMutex

	// Properties holds the index properties.
	Properties map[string]Property `json:"properties,omitempty"`
}

// NewMappings returns a initialized Mappings object.
func NewMappings() *Mappings {
	return &Mappings{
		Properties: make(map[string]Property),
	}
}

// Len returns the number of properties stored.
// This function is concurrent safe.
func (t *Mappings) Len() int {
	t.lock.RLock()
	n := len(t.Properties)
	t.lock.RUnlock()
	return n
}

// SetProperty sets/ adds the given property to the mapping.
// This function is concurrent safe.
func (t *Mappings) SetProperty(field string, prop Property) {
	t.lock.Lock()
	t.Properties[field] = prop
	t.lock.Unlock()
}

// GetProperty returns the property by its field name.
// This function is concurrent safe.
func (t *Mappings) GetProperty(field string) (Property, bool) {
	t.lock.RLock()
	prop, ok := t.Properties[field]
	t.lock.RUnlock()
	return prop, ok
}

// ListProperty returns all properties of the mapping.
// This function is concurrent safe.
func (t *Mappings) ListProperty() map[string]Property {
	m := make(map[string]Property)

	t.lock.RLock()
	for k, v := range t.Properties {
		m[k] = v
	}
	t.lock.RUnlock()

	return m
}

// MarshalJSON overwrites the default marshaler.
// This function is concurrent safe.
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

// Property holds the property of a single field.
type Property struct {
	// Type holds the field type.
	// The type field may be empty, i.e., if the field represents an object.
	Type string `json:"type,omitempty"`
	// Properties holds the sub-properties of this property.
	Properties map[string]Property `json:"properties,omitempty"`
	// Fields the same string value to be indexed in multiple ways for different purposes,
	// such as one field for search and a multi-field for sorting and aggregations,
	// or the same string value analyzed by different analyzers.
	Fields map[string]Property `json:"fields,omitempty"`
	// IgnoreAbove prevents indexing of strings longer than the configured value.
	// TODO: Should this be supported by Zinc?
	IgnoreAbove    uint   `json:"ignore_above,omitempty"`
	Analyzer       string `json:"analyzer,omitempty"`
	SearchAnalyzer string `json:"search_analyzer,omitempty"`
	// Format holds the property format.
	Format string `json:"format,omitempty"`
}

// NewProperty returns a new Property object.
func NewProperty(typ string) Property {
	return Property{
		Type:       typ,
		Properties: make(map[string]Property),
		Fields:     make(map[string]Property),
	}
}
