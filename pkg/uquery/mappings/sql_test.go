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

package mappings

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/pkg/config"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func TestSQLMappingsTypes(t *testing.T) {
	enabled := config.Global.EnableSQLMappings
	config.Global.EnableSQLMappings = true
	t.Cleanup(func() { config.Global.EnableSQLMappings = enabled })
	tests := []struct{ sqlType, mappingType string }{
		{"VARCHAR(255)", "keyword"}, {"CHAR(10)", "keyword"},
		{"ENUM('a,b','c')", "keyword"}, {"SET('a','b')", "keyword"},
		{"TEXT", "text"}, {"TINYTEXT", "text"}, {"MEDIUMTEXT", "text"}, {"LONGTEXT", "text"},
		{"TINYINT(1)", "numeric"}, {"SMALLINT", "numeric"}, {"MEDIUMINT", "numeric"},
		{"INT UNSIGNED", "numeric"}, {"INTEGER", "numeric"}, {"BIGINT", "numeric"},
		{"DECIMAL(10,2)", "numeric"}, {"NUMERIC(10,2)", "numeric"},
		{"FLOAT", "numeric"}, {"DOUBLE", "numeric"}, {"REAL", "numeric"},
		{"YEAR", "numeric"}, {"BIT(8)", "numeric"},
		{"BOOL", "bool"}, {"BOOLEAN", "bool"},
		{"DATE", "date"}, {"DATETIME(6)", "date"}, {"TIMESTAMP", "date"},
		{"TIME", "keyword"}, {"JSON", ""},
	}
	for _, tt := range tests {
		t.Run(tt.sqlType, func(t *testing.T) {
			m, err := Request(nil, map[string]interface{}{"sql": "CREATE TABLE t (`value` " + tt.sqlType + ")"})
			if !assert.NoError(t, err) {
				return
			}
			p, ok := m.GetProperty("value")
			assert.Equal(t, tt.mappingType != "", ok)
			assert.Equal(t, tt.mappingType, p.Type)
		})
	}
}

func TestSQLMappingsDDL(t *testing.T) {
	ddl := "/* schema */ CREATE TABLE IF NOT EXISTS `db`.`users` (" +
		"`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, name VARCHAR(100) DEFAULT 'a,b', " +
		"bio TEXT COMMENT 'text, comment', created_at DATETIME(6), birthday DATE, metadata JSON, " +
		"PRIMARY KEY (id), KEY name_idx (name)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"
	properties, err := sqlProperties(ddl)
	if !assert.NoError(t, err) {
		return
	}
	assert.Len(t, properties, 6)
	m, err := Request(nil, map[string]interface{}{"properties": properties})
	if !assert.NoError(t, err) {
		return
	}
	for field, value := range map[string]string{"created_at": "2026-09-07 12:34:56.123456", "birthday": "2026-09-07"} {
		p, ok := m.GetProperty(field)
		assert.True(t, ok)
		_, err := zutils.ParseTime(value, p.Format, p.TimeZone)
		assert.NoError(t, err)
	}
}

func TestSQLMappingsRejectInvalid(t *testing.T) {
	tests := []string{
		"", "SELECT 1", "DROP TABLE t", "ALTER TABLE t ADD id INT",
		"CREATE TABLE t LIKE other", "CREATE TABLE t AS SELECT 1",
		"CREATE TABLE t ()", "CREATE TABLE t (id INT,)",
		"CREATE TABLE t (id INT) garbage", "CREATE TABLE t (id INT); DROP TABLE t",
		"CREATE TABLE t (id INT); CREATE TABLE u (id INT)",
		"CREATE TABLE t (id INT, ID TEXT)", "CREATE TABLE t (id GEOMETRY)",
		"CREATE TABLE t (id BLOB)", "CREATE TABLE t (id UNKNOWN)",
	}
	for _, ddl := range tests {
		t.Run(ddl, func(t *testing.T) {
			properties, err := sqlProperties(ddl)
			assert.Error(t, err)
			assert.Nil(t, properties)
		})
	}
}

func TestSQLMappingsRequest(t *testing.T) {
	enabled := config.Global.EnableSQLMappings
	t.Cleanup(func() { config.Global.EnableSQLMappings = enabled })
	config.Global.EnableSQLMappings = true
	for _, data := range []map[string]interface{}{
		{"sql": nil}, {"sql": 42}, {"sql": "  "},
		{"sql": "CREATE TABLE t (id INT)", "properties": map[string]interface{}{}},
		{"sql": "CREATE TABLE t (id INT)", "unknown": true},
	} {
		_, err := Request(nil, data)
		assert.Error(t, err)
	}
	config.Global.EnableSQLMappings = false
	_, err := Request(nil, map[string]interface{}{"sql": "CREATE TABLE t (id INT)"})
	assert.ErrorContains(t, err, "disabled")
	m, err := Request(nil, map[string]interface{}{"properties": map[string]interface{}{
		"name":    map[string]interface{}{"type": "keyword", "store": true},
		"address": map[string]interface{}{"properties": map[string]interface{}{"city": map[string]interface{}{"type": "keyword"}}},
	}})
	if !assert.NoError(t, err) {
		return
	}
	p, ok := m.GetProperty("name")
	assert.True(t, ok)
	assert.True(t, p.Store)
	_, ok = m.GetProperty("address.city")
	assert.True(t, ok)
}
