package v1

import (
	"time"
)

// ZincQuery is the query object for the zinc index. All search requests should send this struct
type ZincQuery struct {
	// SearchType is the type of search to perform. Can be match, match_phrase, query_string, etc
	SearchType   string                       `json:"search_type"`
	MaxResults   int                          `json:"max_results"`
	From         int                          `json:"from"`
	Explain      bool                         `json:"explain"`
	Highlight    QueryHighlight               `json:"highlight"`
	Query        QueryParams                  `json:"query"`
	Aggregations map[string]AggregationParams `json:"aggs"`
	SortFields   []string                     `json:"sort_fields"`
	Source       interface{}                  `json:"_source"`
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

// SearchResponse for a query
type SearchResponse struct {
	Took         int                            `json:"took"` // Time it took to generate the response
	TimedOut     bool                           `json:"timed_out"`
	Hits         Hits                           `json:"hits"`
	Aggregations map[string]AggregationResponse `json:"aggregations,omitempty"`
	Error        string                         `json:"error"`
}

type Hits struct {
	Total    Total   `json:"total"`
	MaxScore float64 `json:"max_score"`
	Hits     []Hit   `json:"hits"`
}

type Hit struct {
	Index     string      `json:"_index"`
	Type      string      `json:"_type"`
	ID        string      `json:"_id"`
	Score     float64     `json:"_score"`
	Timestamp time.Time   `json:"@timestamp"`
	Source    interface{} `json:"_source"`
}

type Total struct {
	Value int `json:"value"` // Count of documents returned
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
