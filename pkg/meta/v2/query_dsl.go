package v2

// ZincQuery is the query object for the zinc index. compatible ES Query DSL
type ZincQuery struct {
	Query          map[string]interface{}  `json:"query"`
	Aggregations   map[string]Aggregations `json:"aggs"`
	Highlight      *Highlight              `json:"highlight"`
	Fields         interface{}             `json:"fields"`  // ["field1", "field2.*", {"field": "fieldName", "format": "epoch_millis"}]
	Source         interface{}             `json:"_source"` // true, false, ["field1", "field2.*"]
	Sort           interface{}             `json:"sort"`    // "_sorce", ["+Year","-Year", {"Year": "desc"}, "Date": {"order": "asc"", "format": "yyyy-MM-dd"}}"}]
	Explain        bool                    `json:"explain"`
	From           int                     `json:"from"`
	Size           int                     `json:"size"`
	Timeout        int                     `json:"timeout"`
	TrackTotalHits bool                    `json:"track_total_hits"`
}

type Query struct {
	Bool              *BoolQuery                `json:"bool"`                // .
	Boosting          *BoostingQuery            `json:"boosting"`            // TODO: not implemented
	Match             map[string]interface{}    `json:"match"`               // simple, MatchQuery
	MatchBoolPrefix   map[string]interface{}    `json:"match_bool_prefix"`   // simple, MatchBoolPrefixQuery
	MatchPhrase       map[string]interface{}    `json:"match_phrase"`        // simple, MatchPhraseQuery
	MatchPhrasePrefix map[string]interface{}    `json:"match_phrase_prefix"` // simple, MatchPhrasePrefixQuery
	MultiMatch        *MultiMatchQuery          `json:"multi_match"`         // .
	MatchAll          interface{}               `json:"match_all"`           // just set or null
	MatchNone         interface{}               `json:"match_none"`          // just set or null
	CombinedFields    *CombinedFieldsQuery      `json:"combined_fields"`     // TODO: not implemented
	QueryString       *QueryStringQuery         `json:"query_string"`        // .
	SimpleQueryString *SimpleQueryStringQuery   `json:"simple_query_string"` // .
	Exists            *ExistsQuery              `json:"exists"`              // .
	Ids               *IdsQuery                 `json:"ids"`                 // .
	Range             map[string]*RangeQuery    `json:"range"`               // .
	Fuzzy             map[string]interface{}    `json:"fuzzy"`               // simple, FuzzyQuery
	Prefix            map[string]interface{}    `json:"prefix"`              // simple, PrefixQuery
	Wildcard          map[string]interface{}    `json:"wildcard"`            // simple, WildcardQuery
	Term              map[string]interface{}    `json:"term"`                // simple, TermQuery
	Terms             map[string]interface{}    `json:"terms"`               // .
	TermsSet          map[string]*TermsSetQuery `json:"terms_set"`           // TODO: not implemented
	GeoBoundingBox    interface{}               `json:"geo_bounding_box"`    // TODO: not implemented
	GeoDistance       interface{}               `json:"geo_distance"`        // TODO: not implemented
	GeoPolygon        interface{}               `json:"geo_polygon"`         // TODO: not implemented
	GeoShape          interface{}               `json:"geo_shape"`           // TODO: not implemented
}

type BoolQuery struct {
	Should             interface{} `json:"should"`               // query, [query1, query2]
	Must               interface{} `json:"must"`                 // query, [query1, query2]
	MustNot            interface{} `json:"must_not"`             // query, [query1, query2]
	Filter             interface{} `json:"filter"`               // query, [query1, query2]
	MinimumShouldMatch float64     `json:"minimum_should_match"` // only for should
}

type BoostingQuery struct {
	Positive      interface{} `json:"positive"` // singe or multiple queries
	Negative      interface{} `json:"negative"` // singe or multiple queries
	NegativeBoost float64     `json:"negative_boost"`
}

type MatchQuery struct {
	Query        string      `json:"query"`
	Analyzer     string      `json:"analyzer"`
	Operator     string      `json:"operator"`  // or(default), and
	Fuzziness    interface{} `json:"fuzziness"` // auto, 1,2,3,n
	PrefixLength float64     `json:"prefix_length"`
	Boost        float64     `json:"boost"`
}

type MatchBoolPrefixQuery struct {
	Query    string  `json:"query"`
	Analyzer string  `json:"analyzer"`
	Boost    float64 `json:"boost"`
}

type MatchPhraseQuery struct {
	Query    string  `json:"query"`
	Analyzer string  `json:"analyzer"`
	Boost    float64 `json:"boost"`
}

type MatchPhrasePrefixQuery struct {
	Query    string  `json:"query"`
	Analyzer string  `json:"analyzer"`
	Boost    float64 `json:"boost"`
}

type MultiMatchQuery struct {
	Query              string   `json:"query"`
	Analyzer           string   `json:"analyzer"`
	Fields             []string `json:"fields"`
	Boost              float64  `json:"boost"`
	Type               string   `json:"type"`     // best_fields(default), most_fields, cross_fields, phrase, phrase_prefix, bool_prefix
	Operator           string   `json:"operator"` // or(default), and
	MinimumShouldMatch float64  `json:"minimum_should_match"`
}

type CombinedFieldsQuery struct {
	Query              string   `json:"query"`
	Analyzer           string   `json:"analyzer"`
	Fields             []string `json:"fields"`
	Operator           string   `json:"operator"` // or(default), and
	MinimumShouldMatch float64  `json:"minimum_should_match"`
}

type QueryStringQuery struct {
	Query           string   `json:"query"`
	Analyzer        string   `json:"analyzer"`
	Fields          []string `json:"fields"`
	DefaultField    string   `json:"default_field"`
	DefaultOperator string   `json:"default_operator"` // or(default), and
	Boost           float64  `json:"boost"`
}

type SimpleQueryStringQuery struct {
	Query           string   `json:"query"`
	Analyzer        string   `json:"analyzer"`
	Fields          []string `json:"fields"`
	DefaultOperator string   `json:"default_operator"` // or(default), and
	AllFields       bool     `json:"all_fields"`
	Boost           float64  `json:"boost"`
}

// ExistsQuery
// {"exists":{"field":"field_name"}}
type ExistsQuery struct {
	Field string `json:"field"`
}

// IdsQuery
// {"ids":{"values":["1","2","3"]}}
type IdsQuery struct {
	Values []string `json:"values"`
}

// RangeQuery
// {"range":{"field":{"gte":10,"lte":20}}}
type RangeQuery struct {
	GT       interface{} `json:"gt"`        // null, float64
	GTE      interface{} `json:"gte"`       // null, float64
	LT       interface{} `json:"lt"`        // null, float64
	LTE      interface{} `json:"lte"`       // null, float64
	Format   string      `json:"format"`    // Date format used to convert date values in the query.
	TimeZone string      `json:"time_zone"` // used to convert date values in the query to UTC.
	Boost    float64     `json:"boost"`
}

// RegexpQuery
// {"regexp":{"field":{"value":"[0-9]*"}}}
type RegexpQuery struct {
	Value string  `json:"value"`
	Flags string  `json:"flags"`
	Boost float64 `json:"boost"`
}

// FuzzyQuery
// {"fuzzy":{"field":"value"}}
// {"fuzzy":{"field":{"value":"value","fuzziness":"auto"}}}
type FuzzyQuery struct {
	Value        string      `json:"value"`
	Fuzziness    interface{} `json:"fuzziness"` // auto, 1,2,3,n
	PrefixLength float64     `json:"prefix_length"`
	Boost        float64     `json:"boost"`
}

// PrefixQuery
// {"prefix":{"field":"value"}}
// {"prefix":{"field":{"value":"value","boost":1.0}}}
type PrefixQuery struct {
	Value string  `json:"value"` // You can speed up prefix queries using the index_prefixes mapping parameter.
	Boost float64 `json:"boost"`
}

// WildcardQuery
// {"wildcard": {"field": "*query*"}}
// {"wildcard": {"field": {"value": "*query*", "boost": 1.0}}}
type WildcardQuery struct {
	Value string  `json:"value"`
	Boost float64 `json:"boost"`
}

// TermQuery
// {"term":{"field": "value"}}
// {"term":{"field": {"value": "value", "boost": 1.0}}}
type TermQuery struct {
	Value           interface{} `json:"value"`
	Boost           float64     `json:"boost"`
	CaseInsensitive bool        `json:"case_insensitive"`
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
	Field       string  `json:"field"`
	Size        int     `json:"size"`
	Interval    float64 `json:"interval"`
	Offset      float64 `json:"offset"`
	MinDocCount int     `json:"min_doc_count"`
	Keyed       bool    `json:"keyed"`
}

type AggregationDateHistogram struct {
	Field            string `json:"field"`
	Size             int    `json:"size"`
	FixedInterval    string `json:"fixed_interval"`    // ms,s,m,h,d
	CalendarInterval string `json:"calendar_interval"` // minute,hour,day,week,month,quarter,year
	Format           string `json:"format"`            // format key_as_string
	TimeZone         string `json:"time_zone"`         // time_zone
	MinDocCount      int    `json:"min_doc_count"`
	Keyed            bool   `json:"keyed"`
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
