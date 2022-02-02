package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prabhatsharma/zinc/pkg/routes"
)

var (
	username  = "admin"
	password  = "Complexpass#123"
	indexName = "games3"
	indexData = `{
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

var r *gin.Engine
var once sync.Once

func server() *gin.Engine {
	if r == nil {
		once.Do(func() {
			godotenv.Load()
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
