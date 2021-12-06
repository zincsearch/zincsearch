package main_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

var (
	username  = "admin"
	password  = "Complexpass#123"
	apiHost   = "http://localhost:4080" // change to your api host
	index     = "games3"
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
	queryData1 = `{
	"search_type": "match",
	"query":
	{
		"term": "DEMTSCHENKO",
		"start_time": "2021-06-02T14:28:31.894Z",
		"end_time": "2021-12-30T15:28:31.894Z"
	},
	"fields": ["_all"]
}`

	queryData2 = `{
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

func TestMain(m *testing.M) {
	fmt.Println("test begin...")
	m.Run()
	fmt.Println("test done.")
}

func TestApiIndex(t *testing.T) {
	if err := apiIndex(index, indexData); err != nil {
		t.Error(err)
	}
}

func TestApiQuery_match(t *testing.T) {
	if err := apiQuery(index, queryData1); err != nil {
		t.Error(err)
	}
}
func TestApiQuery_querystring(t *testing.T) {
	if err := apiQuery(index, queryData2); err != nil {
		t.Error(err)
	}
}

func apiIndex(index, data string) error {
	buf := bytes.NewReader([]byte(data))
	req, err := http.NewRequest("PUT", apiHost+"/api/"+index+"/document", buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("http code: %d\nindex response: %s\n", resp.StatusCode, body)
	resp.Body.Close()
	return nil
}

func apiQuery(index, params string) error {
	buf := bytes.NewReader([]byte(params))
	req, err := http.NewRequest("POST", apiHost+"/api/"+index+"/_search", buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("http code: %d\nquery response: %s\n", resp.StatusCode, body)
	resp.Body.Close()
	return nil
}
