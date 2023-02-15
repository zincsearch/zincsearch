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

package v1

import (
	"time"

	"github.com/zinclabs/zincsearch/pkg/meta"
)

// ZincQuery is the query object for the zinc index. All search requests should send this struct
type ZincQuery struct {
	// SearchType is the type of search to perform. Can be match, match_phrase, query_string, etc
	SearchType   string                       `json:"search_type"`
	MaxResults   int                          `json:"max_results"`
	From         int                          `json:"from"`
	Explain      bool                         `json:"explain"`
	Highlight    *meta.Highlight              `json:"highlight"`
	Query        QueryParams                  `json:"query"`
	Aggregations map[string]AggregationParams `json:"aggs"`
	SortFields   []string                     `json:"sort_fields"`
	Source       interface{}                  `json:"_source"`
}

type ZincQueryForSDK struct {
	// SearchType is the type of search to perform. Can be match, match_phrase, query_string, etc
	SearchType   string                       `json:"search_type"`
	MaxResults   int                          `json:"max_results"`
	From         int                          `json:"from"`
	Explain      bool                         `json:"explain"`
	Highlight    *meta.Highlight              `json:"highlight"`
	Query        QueryParams                  `json:"query"`
	Aggregations map[string]AggregationParams `json:"aggs"`
	SortFields   []string                     `json:"sort_fields"`
	Source       []string                     `json:"_source"`
}

type QueryParams struct {
	Boost     int        `json:"boost"`
	Term      string     `json:"term"`
	Terms     [][]string `json:"terms"` // For multi phrase query
	Field     string     `json:"field"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
}

type AggregationParams struct {
	AggType      string                       `json:"agg_type"`
	Field        string                       `json:"field"`
	WeightField  string                       `json:"weight_field"` // Field name to be used for setting weight for primary field for weighted average aggregation
	Size         int                          `json:"size"`
	Ranges       []AggregationNumberRange     `json:"ranges"`
	DateRanges   []AggregationDateRange       `json:"date_ranges"`
	Aggregations map[string]AggregationParams `json:"aggs"`
}

type AggregationNumberRange struct {
	From float64 `json:"from"`
	To   float64 `json:"to"`
}

type AggregationDateRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type QueryHighlight struct {
	Fields []string `json:"fields"`
	Style  string   `json:"style"`
}

type Source struct {
	Enable bool            // enable _source returns, default is true
	Fields map[string]bool // what fields can returns
}

type AggregationResponse struct {
	Value   interface{}         `json:"value,omitempty"`
	Buckets []AggregationBucket `json:"buckets,omitempty"`
}

type AggregationBucket struct {
	Key          string                         `json:"key"`
	DocCount     uint64                         `json:"doc_count"`
	Aggregations map[string]AggregationResponse `json:"aggregations,omitempty"`
}

// SearchResponse for a query
type SearchResponse struct {
	Took         int                            `json:"took"` // Time it took to generate the response
	TimedOut     bool                           `json:"timed_out"`
	Hits         Hits                           `json:"hits"`
	Aggregations map[string]AggregationResponse `json:"aggregations,omitempty"`
	Error        string                         `json:"error,omitempty"`
}

type Hits struct {
	Total    Total   `json:"total"`
	MaxScore float64 `json:"max_score"`
	Hits     []Hit   `json:"hits"`
}

type Hit struct {
	Index     string                 `json:"_index"`
	Type      string                 `json:"_type"`
	ID        string                 `json:"_id"`
	Score     float64                `json:"_score"`
	Timestamp time.Time              `json:"@timestamp"`
	Source    interface{}            `json:"_source"`
	Highlight map[string]interface{} `json:"highlight,omitempty"`
}

type Total struct {
	Value int `json:"value"` // Count of documents returned
}
