package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAnalyze(t *testing.T) {
	Convey("test analyzer", t, func() {
		Convey("standard analyzer", func() {
			input := `{
				"analyzer": "standard",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[the 2 quick brown foxes jumped over the lazy dog's bone]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("standard analyzer with stopwords", func() {
			index := `{
				"settings": {
				  "analysis": {
					"analyzer": {
					  "my_english_analyzer": {
						"type": "standard",
						"stopwords": ["_english_"]
					  }
					}
				  }
				}
			  }`
			input := `{
				"analyzer": "my_english_analyzer",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[2 quick brown foxes jumped lazy dog's bone]`

			// create index with custom analyzer
			body := bytes.NewBuffer(nil)
			body.WriteString(index)
			resp := request("PUT", "/api/index/my-index-001", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			// analyze
			body.Reset()
			body.WriteString(input)
			resp = request("POST", "/api/my-index-001/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)
			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)

			// delete index
			request("DELETE", "/api/index/my-index-001", nil)
		})
	})

	Convey("test tokenizer", t, func() {

	})

	Convey("test char_filter", t, func() {

	})

	Convey("test token_filter", t, func() {

	})
}

func getTokenStrings(data []byte) (string, error) {
	var ret map[string]interface{}
	err := json.Unmarshal(data, &ret)
	if err != nil {
		return "", err
	}

	tokens, _ := ret["tokens"].([]interface{})
	if tokens == nil {
		return "", fmt.Errorf("tokens not exists")
	}

	strs := make([]string, 0, len(tokens))
	for _, token := range tokens {
		str := token.(map[string]interface{})["token"].(string)
		strs = append(strs, str)
	}

	return "[" + strings.Join(strs, " ") + "]", nil
}
