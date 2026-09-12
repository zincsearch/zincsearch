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
	"fmt"
	"strings"

	"vitess.io/vitess/go/vt/sqlparser"
)

// sqlProperties translates schema only; SQL is never executed.
func sqlProperties(ddl string) (map[string]interface{}, error) {
	parser, err := sqlparser.New(sqlparser.Options{})
	if err != nil {
		return nil, fmt.Errorf("initialize SQL parser: %w", err)
	}
	stmt, err := parser.ParseStrictDDL(ddl)
	if err != nil {
		return nil, fmt.Errorf("invalid MySQL CREATE TABLE: %w", err)
	}
	create, ok := stmt.(*sqlparser.CreateTable)
	if !ok || create.TableSpec == nil || len(create.TableSpec.Columns) == 0 {
		return nil, fmt.Errorf("sql must contain one CREATE TABLE with explicit columns")
	}

	properties := make(map[string]interface{}, len(create.TableSpec.Columns))
	seen := make(map[string]bool, len(create.TableSpec.Columns))
	for _, column := range create.TableSpec.Columns {
		name := column.Name.String()
		key := strings.ToLower(name)
		if name == "" || seen[key] {
			return nil, fmt.Errorf("empty or duplicate SQL column [%s]", name)
		}
		seen[key] = true
		prop := make(map[string]interface{})
		switch typ := strings.ToLower(column.Type.Type); typ {
		case "char", "varchar", "enum", "set", "time":
			prop["type"] = "keyword"
		case "tinytext", "text", "mediumtext", "longtext":
			prop["type"] = "text"
		case "tinyint", "smallint", "mediumint", "int", "integer", "bigint", "decimal", "numeric", "float", "double", "real", "year", "bit":
			prop["type"] = "numeric"
		case "bool", "boolean":
			prop["type"] = "bool"
		case "date":
			prop["type"] = "date"
			prop["format"] = "2006-01-02"
		case "datetime", "timestamp":
			prop["type"] = "date"
			prop["format"] = "2006-01-02 15:04:05"
		case "json":
			// Like JSON object mappings, child fields are inferred during ingestion.
			prop["type"] = "object"
		default:
			return nil, fmt.Errorf("SQL column [%s] has unsupported type [%s]", name, typ)
		}
		properties[name] = prop
	}
	return properties, nil
}
