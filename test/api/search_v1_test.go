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
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/zutils/json"
)

func TestSearchV1(t *testing.T) {
	t.Run("init data for search", func(t *testing.T) {
		body := bytes.NewBuffer(nil)
		body.WriteString(indexData)
		resp := request("PUT", "/api/"+indexName+"/_doc", body)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("POST /api/:target/_search", func(t *testing.T) {
		t.Run("search document with not exist indexName", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{}`)
			resp := request("POST", "/api/notExistSearch/_search", body)
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})
		t.Run("search document with exist indexName", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "alldocuments"}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
		t.Run("search document with not exist term", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "match", "query": {"term": "xxxx"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.Equal(t, 0, data.Hits.Total.Value)
		})
		t.Run("search document with exist term", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "match", "query": {"term": "DEMTSCHENKO"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, data.Hits.Total.Value, 1)
		})
		t.Run("search document type: alldocuments", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "alldocuments", "query": {}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, data.Hits.Total.Value, 1)
		})
		t.Run("search document type: wildcard", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "wildcard", "query": {"term": "dem*"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, data.Hits.Total.Value, 1)
		})
		t.Run("search document type: fuzzy", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "fuzzy", "query": {"term": "demtschenk"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, data.Hits.Total.Value, 1)
		})
		t.Run("search document type: term", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{
				"search_type": "term", 
				"query": {
					"term": "turin", 
					"field":"City"
				}
			}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, data.Hits.Total.Value, 1)
		})
		t.Run("search document type: daterange", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(fmt.Sprintf(`{
				"search_type": "daterange",
				"query": {
					"start_time": "%s",
					"end_time": "%s"
				}
			}`,
				time.Now().UTC().Add(time.Hour*-24).Format("2006-01-02T15:04:05Z"),
				time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			))
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, data.Hits.Total.Value, 1)
		})
		t.Run("search document type: matchall", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "matchall", "query": {"term": "demtschenk"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, data.Hits.Total.Value, 1)
		})
		t.Run("search document type: match", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "match", "query": {"term": "DEMTSCHENKO"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, data.Hits.Total.Value, 1)
		})
		t.Run("search document type: matchphrase", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "matchphrase", "query": {"term": "DEMTSCHENKO"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, data.Hits.Total.Value, 1)
		})
		t.Run("search document type: multiphrase", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{
				"search_type": "multiphrase",
				"query": {
					"terms": [
						["demtschenko"],
						["albert"]
					]
				}
			}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, data.Hits.Total.Value, 1)
		})
		t.Run("search document type: prefix", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "prefix", "query": {"term": "dem"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, data.Hits.Total.Value, 1)
		})
		t.Run("search document type: querystring", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{"search_type": "querystring", "query": {"term": "DEMTSCHENKO"}}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, data.Hits.Total.Value, 1)
		})
	})

	t.Run("POST /api/:target/_search with aggregations", func(t *testing.T) {
		t.Run("terms aggregation", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{
				"search_type": "matchall", 
				"aggs": {
					"my-agg": {
						"agg_type": "terms",
						"field": "City"
					}
				}
			}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(data.Aggregations), 1)
		})

		t.Run("metric aggregation", func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			body.WriteString(`{
				"search_type": "matchall", 
				"aggs": {
					"my-agg-max": {
						"agg_type": "max",
						"field": "Year"
					},
					"my-agg-min": {
						"agg_type": "min",
						"field": "Year"
					},
					"my-agg-avg": {
						"agg_type": "avg",
						"field": "Year"
					}
				}
			}`)
			resp := request("POST", "/api/"+indexName+"/_search", body)
			assert.Equal(t, http.StatusOK, resp.Code)

			data := new(meta.SearchResponse)
			err := json.Unmarshal(resp.Body.Bytes(), data)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(data.Aggregations), 1)
		})
	})
}
