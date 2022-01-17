package v1

import (
	"time"

	"github.com/blugelabs/bluge/search"
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
	AggType string             `json:"agg_type"`
	Field   string             `json:"field"`
	Sort    string             `json:"sort"`
	Size    int                `json:"size"`
	Ranges  []AggregationRange `json:"ranges"`
}

type AggregationRange struct {
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

type QueryHighlight struct {
	Fields []string `json:"fields"`
	Style  string   `json:"style"`
}

// SearchResponse for a query
type SearchResponse struct {
	Took     int              `json:"took"` // Time it took to generate the response
	TimedOut bool             `json:"timed_out"`
	MaxScore float64          `json:"max_score"`
	Hits     Hits             `json:"hits"`
	Buckets  []*search.Bucket `json:"buckets"`
	Error    string           `json:"error"`
}

type Hits struct {
	Total Total `json:"total"`
	Hits  []Hit `json:"hits"`
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
