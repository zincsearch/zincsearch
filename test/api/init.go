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

package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/zincsearch/zincsearch/pkg/routes"
)

var (
	username   = "admin"
	password   = "Complexpass#123"
	indexName  = "games3"
	indexAlias = indexName + "-alias"
	indexData  = `{
	"Athlete": "DEMTSCHENKO, Albert",
	"City": "Turin",
	"Country": "RUS",
	"Discipline": "Luge",
	"Event": "Singles",
	"Gender": "Men",
	"Medal": "Silver",
	"Season": "winter",
	"Sport": "Luge",
	"Year": 2006
}`
	// 	queryDataMatch = `{
	// 	"search_type": "match",
	// 	"query":
	// 	{
	// 		"term": "DEMTSCHENKO",
	// 		"start_time": "2021-06-02T14:28:31.894Z",
	// 		"end_time": "2021-12-30T15:28:31.894Z"
	// 	},
	// 	"fields": ["_all"]
	// }`

	// 	queryDataQueryString = `{
	//     "search_type": "querystring",
	//     "query":
	//     {
	//         "term": "+City:Turin +Silver",
	//         "start_time": "2021-06-02T14:28:31.894Z",
	//         "end_time": "2021-12-30T15:28:31.894Z"
	//     },
	//     "fields": ["_all"]
	// }`
	bulkData = `{"index": {"_index": "games3", "_type": "doc", "_id": "1"}}
{"field1": "value1", "field2": "value2", "field3": "value3"}
{"index": {"_index": "games3", "_type": "doc", "_id": "2"}}
{"field1": "value1", "field2": "value2", "field3": "value3"}
{"create": {"_index": "games3", "_type": "doc", "_id": "3"}}
{"field1": "value1", "field2": "value2", "field3": "value3"}`
	bulkDataWithDelete = `{"index": {"_index": "games3", "_type": "doc", "_id": "4"}}
{"field1": "value1", "field2": "value2", "field3": "value3"}
{"delete": {"_index": "games3", "_type": "doc", "_id": "1"}}`
)

var (
	r    *gin.Engine
	once sync.Once
)

func server() *gin.Engine {
	if r == nil {
		once.Do(func() {
			gin.SetMode(gin.ReleaseMode)
			r = gin.New()
			r.Use(gin.Recovery())
			routes.SetRoutes(r)
		})
	}

	return r
}

func request(method, api string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, api, body)
	req.SetBasicAuth(username, password)
	w := httptest.NewRecorder()
	server().ServeHTTP(w, req)
	return w
}
