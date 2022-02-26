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

		Convey("standard analyzer with stopwords and filters", func() {
			index := `{
				"settings": {
				  "analysis": {
					"analyzer": {
					  "my_english_analyzer": {
						"type": "standard",
						"stopwords": ["_english_"],
						"token_filter": ["lowercase", "apostrophe", "my_length"]
					  }
					},
					"token_filter": {
						"my_length": {
							"type": "length",
							"min": 2,
							"max": 10
						}
					}
				  }
				}
			  }`
			input := `{
				"analyzer": "my_english_analyzer",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[quick brown foxes jumped lazy dog bone]`

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

		Convey("simple analyzer", func() {
			input := `{
				"analyzer": "simple",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[the quick brown foxes jumped over the lazy dog s bone]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("keyword analyzer", func() {
			input := `{
				"analyzer": "keyword",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[The 2 QUICK Brown-Foxes jumped over the lazy dog's bone.]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("regexp analyzer", func() {
			input := `{
				"analyzer": "regexp",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[the 2 quick brown foxes jumped over the lazy dog s bone]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("regexp analyzer with pattern", func() {
			index := `{
				"settings": {
				  "analysis": {
					"analyzer": {
					  "my_email_analyzer": {
						"type":      "pattern",
						"pattern":   "[^\\W_]+", 
						"lowercase": true
					  }
					}
				  }
				}
			  }`
			input := `{
				"analyzer": "my_email_analyzer",
				"text": "John_Smith@foo-bar.com"
			  }`
			output := `[john smith foo bar com]`

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

		Convey("stop analyzer", func() {
			input := `{
				"analyzer": "stop",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[quick brown foxes jumped lazy dog s bone]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("stop analyzer with stopwords", func() {
			index := `{
				"settings": {
				  "analysis": {
					"analyzer": {
					  "my_stop_analyzer": {
						"type": "stop",
						"stopwords": ["the", "over"]
					  }
					}
				  }
				}
			  }`
			input := `{
				"analyzer": "my_stop_analyzer",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[quick brown foxes jumped lazy dog s bone]`

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

		Convey("whitespace analyzer", func() {
			input := `{
				"analyzer": "whitespace",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[The 2 QUICK Brown-Foxes jumped over the lazy dog's bone.]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})
	})

	Convey("test tokenizer", t, func() {
		Convey("standard tokenizer", func() {
			input := `{
				"tokenizer": "standard",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[The 2 QUICK Brown Foxes jumped over the lazy dog's bone]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("letter tokenizer", func() {
			input := `{
				"tokenizer": "letter",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[The QUICK Brown Foxes jumped over the lazy dog s bone]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("lowercase tokenizer", func() {
			input := `{
				"tokenizer": "lowercase",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[the quick brown foxes jumped over the lazy dog s bone]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("whitespace tokenizer", func() {
			input := `{
				"tokenizer": "whitespace",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[The 2 QUICK Brown-Foxes jumped over the lazy dog's bone.]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("ngram tokenizer", func() {
			input := `{
				"tokenizer": "ngram",
				"text": "Quick Fox"
			  }`
			output := `[Q Qu u ui i ic c ck k k     F F Fo o ox x]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("ngram tokenizer with configuration", func() {
			index := `{
				"settings": {
				  "analysis": {
					"analyzer": {
					  "my_analyzer": {
						"tokenizer": "my_tokenizer"
					  }
					},
					"tokenizer": {
					  "my_tokenizer": {
						"type": "ngram",
						"min_gram": 3,
						"max_gram": 3,
						"token_chars": [
    					  "letter",
            			  "digit"
          				]
					  }
					}
				  }
				}
			  }`
			input := `{
				"analyzer": "my_analyzer",
				"text": "2 Quick Foxes."
			  }`
			output := `[Qui uic ick Fox oxe xes]`

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

		Convey("edge_ngram tokenizer", func() {
			input := `{
				"tokenizer": "edge_ngram",
				"text": "Quick Fox"
			  }`
			output := `[Q Qu]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("edge_ngram tokenizer with configuration", func() {
			index := `{
				"settings": {
				  "analysis": {
					"analyzer": {
					  "my_analyzer": {
						"tokenizer": "my_tokenizer"
					  }
					},
					"tokenizer": {
					  "my_tokenizer": {
						"type": "edge_ngram",
						"min_gram": 2,
						"max_gram": 10,
						"token_chars": [
						  "letter",
						  "digit"
						]
					  }
					}
				  }
				}
			  }`
			input := `{
				"analyzer": "my_analyzer",
				"text": "2 Quick Foxes."
			  }`
			output := `[Qu Qui Quic Quick Fo Fox Foxe Foxes]`

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

		Convey("keyword tokenizer", func() {
			input := `{
				"tokenizer": "keyword",
				"text": "New York"
			  }`
			output := `[New York]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("keyword tokenizer with filters", func() {
			input := `{
				"tokenizer": "keyword",
				"token_filter": [ "lowercase" ],
				"text": "john.SMITH@example.COM"
			  }`
			output := `[john.smith@example.com]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("regexp tokenizer", func() {
			input := `{
				"tokenizer": "regexp",
				"text": "The foo_bar_size's default is 5."
			  }`
			output := `[The foo_bar_size s default is 5]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("regexp tokenizer with configuration example1", func() {
			index := `{
				"settings": {
				  "analysis": {
					"analyzer": {
					  "my_analyzer": {
						"tokenizer": "my_tokenizer"
					  }
					},
					"tokenizer": {
					  "my_tokenizer": {
						"type": "pattern",
						"pattern": "[^,]+"
					  }
					}
				  }
				}
			  }`
			input := `{
				"analyzer": "my_analyzer",
				"text": "comma,separated,values"
			  }`
			output := `[comma separated values]`

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

		Convey("regexp tokenizer with configuration example2", func() {
			index := `{
				"settings": {
				  "analysis": {
					"analyzer": {
					  "my_analyzer": {
						"tokenizer": "my_tokenizer"
					  }
					},
					"tokenizer": {
					  "my_tokenizer": {
						"type": "pattern",
						"pattern": "((?:\\\\\"|[^\", ])+)"
					  }
					}
				  }
				}
			  }`
			input := `{
				"analyzer": "my_analyzer",
				"text": "\"value\", \"value with embedded \\\" quote\""
			  }`
			output := `[value value with embedded \" quote]`

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

		Convey("char_group tokenizer", func() {
			input := `{
				"tokenizer": {
				  "type": "char_group",
				  "tokenize_on_chars": [
					"whitespace",
					"-",
					"\n"
				  ]
				},
				"text": "The QUICK brown-fox"
			  }`
			output := `[The QUICK brown fox]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("path_hierarchy tokenizer", func() {
			input := `{
				"tokenizer": "path_hierarchy",
				"text": "/one/two/three"
			  }`
			output := `[/one /one/two /one/two/three]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			So(resp.Code, ShouldEqual, http.StatusOK)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			So(err, ShouldBeNil)
			So(tokens, ShouldEqual, output)
		})

		Convey("path_hierarchy tokenizer with configuration", func() {
			index := `{
				"settings": {
				  "analysis": {
					"analyzer": {
					  "my_analyzer": {
						"tokenizer": "my_tokenizer"
					  }
					},
					"tokenizer": {
					  "my_tokenizer": {
						"type": "path_hierarchy",
						"delimiter": "-",
						"replacement": "/",
						"skip": 2
					  }
					}
				  }
				}
			  }`
			input := `{
				"analyzer": "my_analyzer",
				"text": "one-two-three-four-five"
			  }`
			output := `[/three /three/four /three/four/five]`

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
