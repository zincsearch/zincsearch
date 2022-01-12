package test

import (
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prabhatsharma/zinc/pkg/routes"
)

var (
	Username  = "admin"
	Password  = "Complexpass#123"
	Index     = "games3"
	IndexData = `{
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
	QueryDataMatch = `{
	"search_type": "match",
	"query":
	{
		"term": "DEMTSCHENKO",
		"start_time": "2021-06-02T14:28:31.894Z",
		"end_time": "2021-12-30T15:28:31.894Z"
	},
	"fields": ["_all"]
}`

	QueryDataQueryString = `{
    "search_type": "querystring",
    "query":
    {
        "term": "+City:Turin +Silver",
       "start_time": "2021-06-02T14:28:31.894Z",
        "end_time": "2021-12-30T15:28:31.894Z"
    },
    "fields": ["_all"]
}`
)

var r *gin.Engine
var once sync.Once

func Server() *gin.Engine {
	if r == nil {
		once.Do(func() {
			godotenv.Load()
			r = gin.New()
			r.Use(gin.Recovery())
			routes.SetRoutes(r)
			// set env
			os.Setenv("MAX_RESULTS", "10000")
			os.Setenv("FIRST_ADMIN_USER", Username)
			os.Setenv("FIRST_ADMIN_PASSWORD", Password)
			os.Setenv("DATA_PATH", "./data")
		})
	}

	return r
}
