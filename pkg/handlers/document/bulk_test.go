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

package document

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zincsearch/test/utils"
)

func TestBulk(t *testing.T) {
	type args struct {
		code   int
		data   string
		params map[string]string
		result string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "bulk",
			args: args{
				code: http.StatusOK,
				data: `{ "index" : { "_index" : "document.bulk" } } 
				{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "HAJOS, Alfred", "Country": "HUN", "Gender": "Men", "Event": "100M Freestyle", "Medal": "Gold", "Season": "summer"}
				{ "create" : { "_index" : "document.bulk" } } 
				{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "HERSCHMANN, Otto", "Country": "AUT", "Gender": "Men", "Event": "100M Freestyle", "Medal": "Silver", "Season": "summer"}
				{ "update" : { "_index" : "document.bulk", "_id": "1" } } 
				{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "HERSCHMANN, Otto", "Country": "AUT", "Gender": "Men", "Event": "100M Freestyle", "Medal": "Silver", "Season": "summer"}
				{ "delete" : { "_index" : "document.bulk", "_id": "1" } } `,
				params: map[string]string{"target": "document.bulk"},
				result: "",
			},
		},
		{
			name: "error",
			args: args{
				code: http.StatusOK,
				data: `{ "index" : { "_index" : "document.bulk" } } 
				{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "HAJOS, Alfred", "Country": "HUN", "Gender": "Men", "Event": "100M Freestyle", "Medal": "Gold", "Season": "summer"}
				{ "delete" : { "_index" : "document.bulk", "_id": "1"x } } `,
				params: map[string]string{"target": "document.bulk"},
				result: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			utils.SetGinRequestData(c, tt.args.data)
			utils.SetGinRequestParams(c, tt.args.params)
			Bulk(c)
			assert.Equal(t, tt.args.code, w.Code)
			assert.Contains(t, w.Body.String(), tt.args.result)
		})
	}
}

func TestESBulk(t *testing.T) {
	type args struct {
		code   int
		data   string
		params map[string]string
		result string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "bulk",
			args: args{
				code: http.StatusOK,
				data: `{ "index" : { "_index" : "document.esbulk" } } 
				{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "HAJOS, Alfred", "Country": "HUN", "Gender": "Men", "Event": "100M Freestyle", "Medal": "Gold", "Season": "summer"}
				{ "create" : { "_index" : "document.esbulk" } } 
				{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "HERSCHMANN, Otto", "Country": "AUT", "Gender": "Men", "Event": "100M Freestyle", "Medal": "Silver", "Season": "summer"}
				{ "update" : { "_index" : "document.esbulk", "_id": "1" } } 
				{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "HERSCHMANN, Otto", "Country": "AUT", "Gender": "Men", "Event": "100M Freestyle", "Medal": "Silver", "Season": "summer"}
				{ "delete" : { "_index" : "document.esbulk", "_id": "1" } } `,
				params: map[string]string{"target": "document.esbulk"},
				result: "",
			},
		},
		{
			name: "error",
			args: args{
				code: http.StatusOK,
				data: `{ "index" : { "_index" : "document.esbulk" } } 
				{"Year": 1896, "City": "Athens", "Sport": "Aquatics", "Discipline": "Swimming", "Athlete": "HAJOS, Alfred", "Country": "HUN", "Gender": "Men", "Event": "100M Freestyle", "Medal": "Gold", "Season": "summer"}
				{ "delete" : { "_index" : "document.esbulk", "_id": "1"x } } `,
				params: map[string]string{"target": "document.esbulk"},
				result: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			utils.SetGinRequestData(c, tt.args.data)
			utils.SetGinRequestParams(c, tt.args.params)
			ESBulk(c)
			assert.Equal(t, tt.args.code, w.Code)
			assert.Contains(t, w.Body.String(), tt.args.result)
		})
	}
}
