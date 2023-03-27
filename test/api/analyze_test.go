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
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zinclabs/zincsearch/pkg/zutils/json"
)

func TestAnalyze(t *testing.T) {
	t.Run("test analyzer", func(t *testing.T) {
		t.Run("standard analyzer", func(t *testing.T) {
			input := `{
				"analyzer": "standard",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[the 2 quick brown foxes jumped over the lazy dog's bone]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("standard analyzer with stopwords", func(t *testing.T) {
			indexName := "my-index-001"
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
			resp := request("PUT", "/api/index/"+indexName, body)
			buf, err := io.ReadAll(resp.Body)
			fmt.Println(string(buf), err)
			assert.Equal(t, http.StatusOK, resp.Code)

			// analyze
			body.Reset()
			body.WriteString(input)
			resp = request("POST", "/api/"+indexName+"/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)
			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)

			// delete index
			request("DELETE", "/api/index/"+indexName, nil)
		})

		t.Run("standard analyzer with stopwords and filters", func(t *testing.T) {
			indexName := "my-index-002"
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
			resp := request("PUT", "/api/index/"+indexName, body)
			assert.Equal(t, http.StatusOK, resp.Code)

			// analyze
			body.Reset()
			body.WriteString(input)
			resp = request("POST", "/api/"+indexName+"/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)
			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)

			// delete index
			request("DELETE", "/api/index/"+indexName, nil)
		})

		t.Run("simple analyzer", func(t *testing.T) {
			input := `{
				"analyzer": "simple",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[the quick brown foxes jumped over the lazy dog s bone]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("keyword analyzer", func(t *testing.T) {
			input := `{
				"analyzer": "keyword",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[The 2 QUICK Brown-Foxes jumped over the lazy dog's bone.]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("regexp analyzer", func(t *testing.T) {
			input := `{
				"analyzer": "regexp",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[the 2 quick brown foxes jumped over the lazy dog s bone]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("regexp analyzer with pattern", func(t *testing.T) {
			indexName := "my-index-003"
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
			resp := request("PUT", "/api/index/"+indexName, body)
			assert.Equal(t, http.StatusOK, resp.Code)

			// analyze
			body.Reset()
			body.WriteString(input)
			resp = request("POST", "/api/"+indexName+"/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)
			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)

			// delete index
			request("DELETE", "/api/index/"+indexName, nil)
		})

		t.Run("stop analyzer", func(t *testing.T) {
			input := `{
				"analyzer": "stop",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[quick brown foxes jumped lazy dog s bone]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("stop analyzer with stopwords", func(t *testing.T) {
			indexName := "my-index-004"
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
			resp := request("PUT", "/api/index/"+indexName, body)
			assert.Equal(t, http.StatusOK, resp.Code)

			// analyze
			body.Reset()
			body.WriteString(input)
			resp = request("POST", "/api/"+indexName+"/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)
			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)

			// delete index
			request("DELETE", "/api/index/"+indexName, nil)
		})

		t.Run("whitespace analyzer", func(t *testing.T) {
			input := `{
				"analyzer": "whitespace",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[The 2 QUICK Brown-Foxes jumped over the lazy dog's bone.]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("web analyzer", func(t *testing.T) {
			input := `{
				"analyzer": "web",
				"text": "Hello info@blugelabs.com, i come from https://zincsearch-docs.zinc.dev/"
			  }`
			output := `[hello info@blugelabs.com come https://zincsearch-docs.zinc.dev/]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})
	})

	t.Run("test tokenizer", func(t *testing.T) {
		t.Run("standard tokenizer", func(t *testing.T) {
			input := `{
				"tokenizer": "standard",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[The 2 QUICK Brown Foxes jumped over the lazy dog's bone]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("letter tokenizer", func(t *testing.T) {
			input := `{
				"tokenizer": "letter",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[The QUICK Brown Foxes jumped over the lazy dog s bone]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("lowercase tokenizer", func(t *testing.T) {
			input := `{
				"tokenizer": "lowercase",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[the quick brown foxes jumped over the lazy dog s bone]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("whitespace tokenizer", func(t *testing.T) {
			input := `{
				"tokenizer": "whitespace",
				"text": "The 2 QUICK Brown-Foxes jumped over the lazy dog's bone."
			  }`
			output := `[The 2 QUICK Brown-Foxes jumped over the lazy dog's bone.]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("n-gram tokenizer", func(t *testing.T) {
			input := `{
				"tokenizer": "ngram",
				"text": "Quick Fox"
			  }`
			output := `[Q Qu u ui i ic c ck k k     F F Fo o ox x]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("n-gram tokenizer with configuration", func(t *testing.T) {
			indexName := "my-index-005"
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
			resp := request("PUT", "/api/index/"+indexName, body)
			assert.Equal(t, http.StatusOK, resp.Code)

			// analyze
			body.Reset()
			body.WriteString(input)
			resp = request("POST", "/api/"+indexName+"/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)
			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)

			// delete index
			request("DELETE", "/api/index/"+indexName, nil)
		})

		t.Run("edge n-gram tokenizer", func(t *testing.T) {
			input := `{
				"tokenizer": "edge_ngram",
				"text": "Quick Fox"
			  }`
			output := `[Q Qu]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("edge n-gram tokenizer with configuration", func(t *testing.T) {
			indexName := "my-index-006"
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
			resp := request("PUT", "/api/index/"+indexName, body)
			assert.Equal(t, http.StatusOK, resp.Code)

			// analyze
			body.Reset()
			body.WriteString(input)
			resp = request("POST", "/api/"+indexName+"/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)
			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)

			// delete index
			request("DELETE", "/api/index/"+indexName, nil)
		})

		t.Run("keyword tokenizer", func(t *testing.T) {
			input := `{
				"tokenizer": "keyword",
				"text": "New York"
			  }`
			output := `[New York]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("keyword tokenizer with filters", func(t *testing.T) {
			input := `{
				"tokenizer": "keyword",
				"token_filter": [ "lowercase" ],
				"text": "john.SMITH@example.COM"
			  }`
			output := `[john.smith@example.com]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("regexp tokenizer", func(t *testing.T) {
			input := `{
				"tokenizer": "regexp",
				"text": "The foo_bar_size's default is 5."
			  }`
			output := `[The foo_bar_size s default is 5]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("regexp tokenizer with configuration example1", func(t *testing.T) {
			indexName := "my-index-007"
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
			resp := request("PUT", "/api/index/"+indexName, body)
			assert.Equal(t, http.StatusOK, resp.Code)

			// analyze
			body.Reset()
			body.WriteString(input)
			resp = request("POST", "/api/"+indexName+"/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)
			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)

			// delete index
			request("DELETE", "/api/index/"+indexName, nil)
		})

		t.Run("regexp tokenizer with configuration example2", func(t *testing.T) {
			indexName := "my-index-008"
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
			resp := request("PUT", "/api/index/"+indexName, body)
			assert.Equal(t, http.StatusOK, resp.Code)

			// analyze
			body.Reset()
			body.WriteString(input)
			resp = request("POST", "/api/"+indexName+"/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)
			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)

			// delete index
			request("DELETE", "/api/index/"+indexName, nil)
		})

		t.Run("character group tokenizer", func(t *testing.T) {
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
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("path hierarchy tokenizer", func(t *testing.T) {
			input := `{
				"tokenizer": "path_hierarchy",
				"text": "/one/two/three"
			  }`
			output := `[/one /one/two /one/two/three]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("path hierarchy tokenizer with configuration", func(t *testing.T) {
			indexName := "my-index-009"
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
			resp := request("PUT", "/api/index/"+indexName, body)
			assert.Equal(t, http.StatusOK, resp.Code)

			// analyze
			body.Reset()
			body.WriteString(input)
			resp = request("POST", "/api/"+indexName+"/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)
			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)

			// delete index
			request("DELETE", "/api/index/"+indexName, nil)
		})
	})

	t.Run("test token filter", func(t *testing.T) {
		t.Run("Apostrophe token filter", func(t *testing.T) {
			input := `{
				"tokenizer" : "standard",
				"filter" : ["apostrophe"],
				"text" : "Istanbul'a veya Istanbul'dan"
			  }`
			output := `[Istanbul veya Istanbul]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("CJK bigram token filter", func(t *testing.T) {
			input := `{
				"tokenizer" : "standard",
				"filter" : ["cjk_bigram"],
				"text" : "東京都は、日本の首都であり"
			  }`
			output := `[東京 京都 都は 日本 本の の首 首都 都で であ あり]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("CJK width token filter", func(t *testing.T) {
			input := `{
				"tokenizer" : "standard",
				"filter" : ["cjk_width"],
				"text" : "ｼｰｻｲﾄﾞﾗｲﾅｰ"
			  }`
			output := `[シーサイドライナー]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Dictionary token filter", func(t *testing.T) {
			input := `{
				"tokenizer": "standard",
				"filter": [
				  {
					"type": "dict",
					"words": ["Donau", "dampf", "meer", "schiff"]
				  }
				],
				"text": "Donaudampfschiff"
			  }`
			output := `[Donaudampfschiff Donau dampf schiff]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Edge n-gram token filter", func(t *testing.T) {
			input := `{
				"tokenizer": "standard",
				"filter": [
				  { "type": "edge_ngram",
					"min_gram": 1,
					"max_gram": 2
				  }
				],
				"text": "the quick brown fox jumps"
			  }`
			output := `[t th q qu b br f fo j ju]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("N-gram token filter", func(t *testing.T) {
			input := `{
				"tokenizer": "standard",
				"filter": [ "ngram" ],
				"text": "Quick fox"
			  }`
			output := `[Q Qu u ui i ic c ck k f fo o ox x]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Elision token filter", func(t *testing.T) {
			input := `{
				"tokenizer" : "standard",
				"filter" : ["elision"],
				"text" : "j’examine près du wharf"
			  }`
			output := `[examine près du wharf]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Stemmer token filter", func(t *testing.T) {
			input := `{
				"tokenizer": "whitespace",
				"filter": [ "stemmer" ],
				"text": "fox running and jumping"
			  }`
			output := `[fox run and jump]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Keyword token filter", func(t *testing.T) {
			input := `{
				"tokenizer": "whitespace",
				"filter": [
				  {
					"type": "keyword",
					"keywords": [ "jumping" ]
				  },
				  "stemmer"
				],
				"text": "fox running and jumping"
			  }`
			output := `[fox run and jumping]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Length token filter", func(t *testing.T) {
			input := `{
				"tokenizer": "whitespace",
				"filter": [
				  {
					"type": "length",
					"min": 0,
					"max": 4
				  }
				],
				"text": "the quick brown fox jumps over the lazy dog"
			  }`
			output := `[the fox over the lazy dog]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Lowercase token filter", func(t *testing.T) {
			input := `{
				"tokenizer" : "standard",
				"filter" : ["lowercase"],
				"text" : "THE Quick FoX JUMPs"
			  }`
			output := `[the quick fox jumps]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Pattern replace token filter", func(t *testing.T) {
			input := `{
				"tokenizer": "whitespace",
				"filter": [
				  {
					"type": "pattern_replace",
					"pattern": "(dog)",
					"replacement": "watch$1"
				  }
				],
				"text": "foxes jump lazy dogs"
			  }`
			output := `[foxes jump lazy watchdogs]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Reverse token filter", func(t *testing.T) {
			input := `{
				"tokenizer" : "standard",
				"filter" : ["reverse"],
				"text" : "quick fox jumps"
			  }`
			output := `[kciuq xof spmuj]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Shingle token filter", func(t *testing.T) {
			input := `{
				"tokenizer": "whitespace",
				"filter": [ "shingle" ],
				"text": "quick brown fox jumps"
			  }`
			output := `[quick brown quick brown fox brown fox jumps fox jumps]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Stop token filter", func(t *testing.T) {
			input := `{
				"tokenizer": "standard",
				"filter": [ "stop" ],
				"text": "a quick fox jumps over the lazy dog"
			  }`
			output := `[quick fox jumps lazy dog]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Stop token filter with configuration", func(t *testing.T) {
			input := `{
				"tokenizer": "standard",
				"filter": [
					{
						"type": "stop",
						"stopwords": ["a", "the", "dog"]
					}
				],
				"text": "a quick fox jumps over the lazy dog"
			}`
			output := `[quick fox jumps over lazy]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Stop token filter with configuration language", func(t *testing.T) {
			input := `{
				"tokenizer": "standard",
				"filter": [
					{
						"type": "stop",
						"stopwords": ["_english_"]
					}
				],
				"text": "a quick fox jumps over the lazy dog"
			  }`
			output := `[quick fox jumps lazy dog]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Trim token filter", func(t *testing.T) {
			input := `{
				"tokenizer" : "keyword",
				"filter" : ["trim"],
				"text" : " fox "
			  }`
			output := `[fox]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Truncate token filter", func(t *testing.T) {
			input := `{
				"tokenizer" : "whitespace",
				"filter" : ["truncate"],
				"text" : "the quinquennial extravaganza carried on"
			  }`
			output := `[the quinquenni extravagan carried on]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Unique token filter", func(t *testing.T) {
			input := `{
				"tokenizer" : "whitespace",
				"filter" : ["unique"],
				"text" : "the quick fox jumps the lazy fox"
			  }`
			output := `[the quick fox jumps lazy]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Uppercase token filter", func(t *testing.T) {
			input := `{
				"tokenizer" : "standard",
				"filter" : ["uppercase"],
				"text" : "the Quick FoX JUMPs"
			  }`
			output := `[THE QUICK FOX JUMPS]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

	})

	t.Run("test character filter", func(t *testing.T) {
		t.Run("ASCII folding character filter", func(t *testing.T) {
			input := `{
				"tokenizer" : "standard",
				"char_filter" : ["asciifolding"],
				"text" : "açaí à la carte"
			  }`
			output := `[acai a la carte]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("HTML strip character filter", func(t *testing.T) {
			input := `{
				"tokenizer" : "standard",
				"char_filter" : ["html_strip"],
				"text": "<p>I'm so <b>happy</b>!</p>"
			  }`
			output := `[I'm so happy]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Mapping character filter", func(t *testing.T) {
			input := `{
				"tokenizer": "keyword",
				"char_filter": [
				  {
					"type": "mapping",
					"mappings": [
					  "٠ => 0",
					  "١ => 1",
					  "٢ => 2",
					  "٣ => 3",
					  "٤ => 4",
					  "٥ => 5",
					  "٦ => 6",
					  "٧ => 7",
					  "٨ => 8",
					  "٩ => 9"
					]
				  }
				],
				"text": "My license plate is ٢٥٠١٥"
			  }`
			output := `[My license plate is 25015]`

			body := bytes.NewBuffer(nil)
			body.WriteString(input)
			resp := request("POST", "/api/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)
		})

		t.Run("Pattern replace character filter", func(t *testing.T) {
			indexName := "my-index-010"
			index := `{
				"settings": {
				  "analysis": {
					"analyzer": {
					  "my_analyzer": {
						"tokenizer": "standard",
						"char_filter": [
						  "my_char_filter"
						]
					  }
					},
					"char_filter": {
					  "my_char_filter": {
						"type": "pattern_replace",
						"pattern": "(\\d+)-",
						"replacement": "${1}_"
					  }
					}
				  }
				}
			  }`
			input := `{
				"analyzer": "my_analyzer",
				"text": "My credit card is 123-456-789"
			  }`
			output := `[My credit card is 123_456_789]`

			// create index with custom analyzer
			body := bytes.NewBuffer(nil)
			body.WriteString(index)
			resp := request("PUT", "/api/index/"+indexName, body)
			assert.Equal(t, http.StatusOK, resp.Code)

			// analyze
			body.Reset()
			body.WriteString(input)
			resp = request("POST", "/api/"+indexName+"/_analyze", body)
			assert.Equal(t, http.StatusOK, resp.Code)
			tokens, err := getTokenStrings(resp.Body.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, output, tokens)

			// delete index
			request("DELETE", "/api/index/"+indexName, nil)
		})
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
