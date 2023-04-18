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

package meta

import "github.com/zincsearch/zincsearch/pkg/bluge/aggregation"

// ZincQuery is the query object for the zinc index. compatible ES Query DSL
type ZincQuery struct {
	Query          interface{}             `json:"query"`
	Aggregations   map[string]Aggregations `json:"aggs"`
	Highlight      *Highlight              `json:"highlight"`
	Fields         interface{}             `json:"fields"`  // ["field1", "field2.*", {"field": "fieldName", "format": "epoch_millis"}]
	Source         interface{}             `json:"_source"` // true, false, ["field1", "field2.*"]
	Sort           interface{}             `json:"sort"`    // "_score", ["+Year","-Year", {"Year": "desc"}, "Date": {"order": "asc"", "format": "yyyy-MM-dd"}}"}]
	Explain        bool                    `json:"explain"`
	From           int                     `json:"from"`
	Size           int                     `json:"size"`
	Timeout        int                     `json:"timeout"`
	TrackTotalHits bool                    `json:"track_total_hits"`
}

type ZincQueryForSDK struct {
	Query          QueryForSDK             `json:"query"`
	Aggregations   map[string]Aggregations `json:"aggs"`
	Highlight      *Highlight              `json:"highlight"`
	Fields         []string                `json:"fields"`  // ["field1", "field2.*", {"field": "fieldName", "format": "epoch_millis"}]
	Source         []string                `json:"_source"` // true, false, ["field1", "field2.*"]
	Sort           []string                `json:"sort"`    // "_score", ["+Year","-Year", {"Year": "desc"}, "Date": {"order": "asc"", "format": "yyyy-MM-dd"}}"}]
	Explain        bool                    `json:"explain"`
	From           int                     `json:"from"`
	Size           int                     `json:"size"`
	Timeout        int                     `json:"timeout"`
	TrackTotalHits bool                    `json:"track_total_hits"`
}

type Query struct {
	Bool              *BoolQuery                         `json:"bool,omitempty"`                // .
	Boosting          *BoostingQuery                     `json:"boosting,omitempty"`            // TODO: not implemented
	Match             map[string]*MatchQuery             `json:"match,omitempty"`               // simple, MatchQuery
	MatchBoolPrefix   map[string]*MatchBoolPrefixQuery   `json:"match_bool_prefix,omitempty"`   // simple, MatchBoolPrefixQuery
	MatchPhrase       map[string]*MatchPhraseQuery       `json:"match_phrase,omitempty"`        // simple, MatchPhraseQuery
	MatchPhrasePrefix map[string]*MatchPhrasePrefixQuery `json:"match_phrase_prefix,omitempty"` // simple, MatchPhrasePrefixQuery
	MultiMatch        *MultiMatchQuery                   `json:"multi_match,omitempty"`         // .
	MatchAll          *MatchAllQuery                     `json:"match_all,omitempty"`           // just set or null
	MatchNone         *MatchNoneQuery                    `json:"match_none,omitempty"`          // just set or null
	CombinedFields    *CombinedFieldsQuery               `json:"combined_fields,omitempty"`     // TODO: not implemented
	QueryString       *QueryStringQuery                  `json:"query_string,omitempty"`        // .
	SimpleQueryString *SimpleQueryStringQuery            `json:"simple_query_string,omitempty"` // .
	Exists            *ExistsQuery                       `json:"exists,omitempty"`              // .
	Ids               *IdsQuery                          `json:"ids,omitempty"`                 // .
	Range             map[string]*RangeQuery             `json:"range,omitempty"`               // simple, FuzzyQuery
	Regexp            map[string]*RegexpQuery            `json:"regexp,omitempty"`              // simple, FuzzyQuery
	Prefix            map[string]*PrefixQuery            `json:"prefix,omitempty"`              // .
	Fuzzy             map[string]*FuzzyQuery             `json:"fuzzy,omitempty"`               // simple, PrefixQuery
	Wildcard          map[string]*WildcardQuery          `json:"wildcard,omitempty"`            // simple, WildcardQuery
	Term              map[string]*TermQuery              `json:"term,omitempty"`                // simple, TermQuery
	Terms             map[string]*TermsQuery             `json:"terms,omitempty"`               // .
	TermsSet          map[string]*TermsSetQuery          `json:"terms_set,omitempty"`           // TODO: not implemented
	GeoBoundingBox    interface{}                        `json:"geo_bounding_box,omitempty"`    // TODO: not implemented
	GeoDistance       interface{}                        `json:"geo_distance,omitempty"`        // TODO: not implemented
	GeoPolygon        interface{}                        `json:"geo_polygon,omitempty"`         // TODO: not implemented
	GeoShape          interface{}                        `json:"geo_shape,omitempty"`           // TODO: not implemented
}

type QueryForSDK struct {
	Bool              *BoolQueryForSDK                   `json:"bool,omitempty"`                // .
	Match             map[string]*MatchQuery             `json:"match,omitempty"`               // simple, MatchQuery
	MatchBoolPrefix   map[string]*MatchBoolPrefixQuery   `json:"match_bool_prefix,omitempty"`   // simple, MatchBoolPrefixQuery
	MatchPhrase       map[string]*MatchPhraseQuery       `json:"match_phrase,omitempty"`        // simple, MatchPhraseQuery
	MatchPhrasePrefix map[string]*MatchPhrasePrefixQuery `json:"match_phrase_prefix,omitempty"` // simple, MatchPhrasePrefixQuery
	MultiMatch        *MultiMatchQuery                   `json:"multi_match,omitempty"`         // .
	MatchAll          *MatchAllQuery                     `json:"match_all,omitempty"`           // just set or null
	MatchNone         *MatchNoneQuery                    `json:"match_none,omitempty"`          // just set or null
	QueryString       *QueryStringQuery                  `json:"query_string,omitempty"`        // .
	SimpleQueryString *SimpleQueryStringQuery            `json:"simple_query_string,omitempty"` // .
	Exists            *ExistsQuery                       `json:"exists,omitempty"`              // .
	Ids               *IdsQuery                          `json:"ids,omitempty"`                 // .
	Range             map[string]*RangeQueryForSDK       `json:"range,omitempty"`               // simple, FuzzyQuery
	Regexp            map[string]*RegexpQuery            `json:"regexp,omitempty"`              // simple, FuzzyQuery
	Prefix            map[string]*PrefixQuery            `json:"prefix,omitempty"`              // .
	Fuzzy             map[string]*FuzzyQuery             `json:"fuzzy,omitempty"`               // simple, PrefixQuery
	Wildcard          map[string]*WildcardQuery          `json:"wildcard,omitempty"`            // simple, WildcardQuery
	Term              map[string]*TermQueryForSDK        `json:"term,omitempty"`                // simple, TermQuery
	Terms             map[string]*TermsQuery             `json:"terms,omitempty"`               // .
}

type BoolQuery struct {
	Should             interface{} `json:"should,omitempty"`               // query, [query1, query2]
	Must               interface{} `json:"must,omitempty"`                 // query, [query1, query2]
	MustNot            interface{} `json:"must_not,omitempty"`             // query, [query1, query2]
	Filter             interface{} `json:"filter,omitempty"`               // query, [query1, query2]
	MinimumShouldMatch float64     `json:"minimum_should_match,omitempty"` // only for should
}

type BoolQueryForSDK struct {
	Should             []*QueryForSDK `json:"should,omitempty"`               // query, [query1, query2]
	Must               []*QueryForSDK `json:"must,omitempty"`                 // query, [query1, query2]
	MustNot            []*QueryForSDK `json:"must_not,omitempty"`             // query, [query1, query2]
	Filter             []*QueryForSDK `json:"filter,omitempty"`               // query, [query1, query2]
	MinimumShouldMatch float64        `json:"minimum_should_match,omitempty"` // only for should
}

type BoostingQuery struct {
	Positive      interface{} `json:"positive,omitempty"` // singe or multiple queries
	Negative      interface{} `json:"negative,omitempty"` // singe or multiple queries
	NegativeBoost float64     `json:"negative_boost,omitempty"`
}

type MatchAllQuery struct{}

type MatchNoneQuery struct{}

type MatchQuery struct {
	Query        string      `json:"query,omitempty"`
	Analyzer     string      `json:"analyzer,omitempty"`
	Operator     string      `json:"operator,omitempty"`  // or(default), and
	Fuzziness    interface{} `json:"fuzziness,omitempty"` // auto, 1,2,3,n
	PrefixLength float64     `json:"prefix_length,omitempty"`
	Boost        float64     `json:"boost,omitempty"`
}

type MatchBoolPrefixQuery struct {
	Query    string  `json:"query,omitempty"`
	Analyzer string  `json:"analyzer,omitempty"`
	Boost    float64 `json:"boost,omitempty"`
}

type MatchPhraseQuery struct {
	Query    string  `json:"query,omitempty"`
	Analyzer string  `json:"analyzer,omitempty"`
	Boost    float64 `json:"boost,omitempty"`
}

type MatchPhrasePrefixQuery struct {
	Query    string  `json:"query,omitempty"`
	Analyzer string  `json:"analyzer,omitempty"`
	Boost    float64 `json:"boost,omitempty"`
}

type MultiMatchQuery struct {
	Query              string   `json:"query,omitempty"`
	Analyzer           string   `json:"analyzer,omitempty"`
	Fields             []string `json:"fields,omitempty"`
	Boost              float64  `json:"boost,omitempty"`
	Type               string   `json:"type,omitempty"`     // best_fields(default), most_fields, cross_fields, phrase, phrase_prefix, bool_prefix
	Operator           string   `json:"operator,omitempty"` // or(default), and
	MinimumShouldMatch float64  `json:"minimum_should_match,omitempty"`
}

type CombinedFieldsQuery struct {
	Query              string   `json:"query,omitempty"`
	Analyzer           string   `json:"analyzer,omitempty"`
	Fields             []string `json:"fields,omitempty"`
	Operator           string   `json:"operator,omitempty"` // or(default), and
	MinimumShouldMatch float64  `json:"minimum_should_match,omitempty"`
}

type QueryStringQuery struct {
	Query           string   `json:"query,omitempty"`
	Analyzer        string   `json:"analyzer,omitempty"`
	Fields          []string `json:"fields,omitempty"`
	DefaultField    string   `json:"default_field,omitempty"`
	DefaultOperator string   `json:"default_operator,omitempty"` // or(default), and
	Boost           float64  `json:"boost,omitempty"`
}

type SimpleQueryStringQuery struct {
	Query           string   `json:"query,omitempty"`
	Analyzer        string   `json:"analyzer,omitempty"`
	Fields          []string `json:"fields,omitempty"`
	DefaultOperator string   `json:"default_operator,omitempty"` // or(default), and
	AllFields       bool     `json:"all_fields,omitempty"`
	Boost           float64  `json:"boost,omitempty"`
}

// ExistsQuery
// {"exists":{"field":"field_name"}}
type ExistsQuery struct {
	Field string `json:"field,omitempty"`
}

// IdsQuery
// {"ids":{"values":["1","2","3"]}}
type IdsQuery struct {
	Values []string `json:"values,omitempty"`
}

// RangeQuery
// {"range":{"field":{"gte":10,"lte":20}}}
type RangeQuery struct {
	GT       interface{} `json:"gt,omitempty"`        // null, float64
	GTE      interface{} `json:"gte,omitempty"`       // null, float64
	LT       interface{} `json:"lt,omitempty"`        // null, float64
	LTE      interface{} `json:"lte,omitempty"`       // null, float64
	Format   string      `json:"format,omitempty"`    // Date format used to convert date values in the query.
	TimeZone string      `json:"time_zone,omitempty"` // used to convert date values in the query to UTC.
	Boost    float64     `json:"boost,omitempty"`
}

type RangeQueryForSDK struct {
	GT       string  `json:"gt,omitempty"`        // string, float64
	GTE      string  `json:"gte,omitempty"`       // string, float64
	LT       string  `json:"lt,omitempty"`        // string, float64
	LTE      string  `json:"lte,omitempty"`       // string, float64
	Format   string  `json:"format,omitempty"`    // Date format used to convert date values in the query.
	TimeZone string  `json:"time_zone,omitempty"` // used to convert date values in the query to UTC.
	Boost    float64 `json:"boost,omitempty"`
}

// RegexpQuery
// {"regexp":{"field":{"value":"[0-9]*"}}}
type RegexpQuery struct {
	Value string  `json:"value,omitempty"`
	Flags string  `json:"flags,omitempty"`
	Boost float64 `json:"boost,omitempty"`
}

// FuzzyQuery
// {"fuzzy":{"field":"value"}}
// {"fuzzy":{"field":{"value":"value","fuzziness":"auto"}}}
type FuzzyQuery struct {
	Value        string      `json:"value,omitempty"`
	Fuzziness    interface{} `json:"fuzziness,omitempty"` // auto, 1,2,3,n
	PrefixLength float64     `json:"prefix_length,omitempty"`
	Boost        float64     `json:"boost,omitempty"`
}

// PrefixQuery
// {"prefix":{"field":"value"}}
// {"prefix":{"field":{"value":"value","boost":1.0}}}
type PrefixQuery struct {
	Value string  `json:"value,omitempty"` // You can speed up prefix queries using the index_prefixes mapping parameter.
	Boost float64 `json:"boost,omitempty"`
}

// WildcardQuery
// {"wildcard": {"field": "*query*"}}
// {"wildcard": {"field": {"value": "*query*", "boost": 1.0}}}
type WildcardQuery struct {
	Value string  `json:"value,omitempty"`
	Boost float64 `json:"boost,omitempty"`
}

// TermQuery
// {"term":{"field": "value"}}
// {"term":{"field": {"value": "value", "boost": 1.0}}}
type TermQuery struct {
	Value           interface{} `json:"value,omitempty"`
	Boost           float64     `json:"boost,omitempty"`
	CaseInsensitive bool        `json:"case_insensitive,omitempty"`
}

type TermQueryForSDK struct {
	Value           string  `json:"value,omitempty"`
	Boost           float64 `json:"boost,omitempty"`
	CaseInsensitive bool    `json:"case_insensitive,omitempty"`
}

// TermsQuery
// {"terms": {"field": ["value1", "value2"], "boost": 1.0}}
type TermsQuery map[string]interface{}

// TermsSetQuery ...
type TermsSetQuery struct{}

type Aggregations struct {
	Avg               *AggregationMetric            `json:"avg"`
	WeightedAvg       *AggregationMetric            `json:"weighted_avg"`
	Max               *AggregationMetric            `json:"max"`
	Min               *AggregationMetric            `json:"min"`
	Sum               *AggregationMetric            `json:"sum"`
	Count             *AggregationMetric            `json:"count"`
	Cardinality       *AggregationMetric            `json:"cardinality"`
	Terms             *AggregationsTerms            `json:"terms"`
	Range             *AggregationRange             `json:"range"`
	DateRange         *AggregationDateRange         `json:"date_range"`
	Histogram         *AggregationHistogram         `json:"histogram"`
	DateHistogram     *AggregationDateHistogram     `json:"date_histogram"`
	AutoDateHistogram *AggregationAutoDateHistogram `json:"auto_date_histogram"`
	IPRange           *AggregationIPRange           `json:"ip_range"` // TODO: not implemented
	Aggregations      map[string]Aggregations       `json:"aggs"`     // nested aggregations
}

type AggregationMetric struct {
	Field       string `json:"field"`
	WeightField string `json:"weight_field"` // Field name to be used for setting weight for primary field for weighted average aggregation
}

type AggregationsTerms struct {
	Field string            `json:"field"`
	Size  int               `json:"size"`
	Order map[string]string `json:"order"` // { "_count": "asc" }
}

type AggregationRange struct {
	Field  string  `json:"field"`
	Ranges []Range `json:"ranges"`
	Keyed  bool    `json:"keyed"`
}

type Range struct {
	To   float64 `json:"to"`
	From float64 `json:"from"`
}

// AggregationDateRange struct
// DateFormat/Pattern refer to:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/search-aggregations-bucket-daterange-aggregation.html#date-format-pattern
type AggregationDateRange struct {
	Field    string      `json:"field"`
	Format   string      `json:"format"`    // format the `to` and `from` values to `_as_string`, and used to format `keyed response`
	TimeZone string      `json:"time_zone"` // refer
	Ranges   []DateRange `json:"ranges"`    // refer
	Keyed    bool        `json:"keyed"`
}

type DateRange struct {
	To   string `json:"to"`
	From string `json:"from"`
}

type AggregationIPRange struct {
	Field  string    `json:"field"`
	Ranges []IPRange `json:"ranges"`
	Keyed  bool      `json:"keyed"`
}

type IPRange struct {
	To   string `json:"to"`
	From string `json:"from"`
}

type AggregationHistogram struct {
	Field          string                      `json:"field"`
	Size           int                         `json:"size"`
	Interval       float64                     `json:"interval"`
	Offset         float64                     `json:"offset"`
	MinDocCount    int                         `json:"min_doc_count"`
	Keyed          bool                        `json:"keyed"`
	ExtendedBounds *aggregation.HistogramBound `json:"extended_bounds"`
	HardBounds     *aggregation.HistogramBound `json:"hard_bounds"`
}

type AggregationDateHistogram struct {
	Field            string                      `json:"field"`
	Size             int                         `json:"size"`
	Interval         string                      `json:"interval"`          // ms,s,m,h,d
	FixedInterval    string                      `json:"fixed_interval"`    // ms,s,m,h,d
	CalendarInterval string                      `json:"calendar_interval"` // minute,hour,day,week,month,quarter,year
	Format           string                      `json:"format"`            // format key_as_string
	TimeZone         string                      `json:"time_zone"`         // time_zone
	MinDocCount      int                         `json:"min_doc_count"`
	Keyed            bool                        `json:"keyed"`
	ExtendedBounds   *aggregation.HistogramBound `json:"extended_bounds"`
	HardBounds       *aggregation.HistogramBound `json:"hard_bounds"`
}

type AggregationAutoDateHistogram struct {
	Field           string `json:"field"`
	Buckets         int    `json:"buckets"`
	MinimumInterval string `json:"minimum_interval"` // minute,hour,day,week,month,quarter,year
	Format          string `json:"format"`           // format key_as_string
	TimeZone        string `json:"time_zone"`        // time_zone
	Keyed           bool   `json:"keyed"`
}

type Highlight struct {
	NumberOfFragments int                   `json:"number_of_fragments"`
	FragmentSize      int                   `json:"fragment_size"`
	PreTags           []string              `json:"pre_tags"`
	PostTags          []string              `json:"post_tags"`
	Fields            map[string]*Highlight `json:"fields"`
}

type Field struct {
	Field  string `json:"field"`
	Format string `json:"format"`
}

type Source struct {
	Enable bool     // enable _source returns, default is true
	Fields []string // what fields can returns
}
