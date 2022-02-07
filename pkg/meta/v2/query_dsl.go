package v2

type ZincQuery struct {
	Query          Query                   `json:"query"`
	Aggregations   map[string]Aggregations `json:"aggs"`
	Highlight      Highlight               `json:"highlight"`
	Fields         []interface{}           `json:"fields"`  // ["field1", "field2.*", {"field": "fieldName", "format": "epoch_millis"}]
	Source         interface{}             `json:"_source"` // true, false, ["field1", "field2.*"]
	Sort           interface{}             `json:"sort"`    // "_sorce", ["Year","-Year", {"Year", "desc"}]
	From           int64                   `json:"from"`
	Size           int64                   `json:"size"`
	Timeout        int64                   `json:"timeout"`
	TrackTotalHits bool                    `json:"track_total_hits"`
}

type Query struct {
	Bool              *BoolQuery              `json:"bool"`                // .
	Boosting          *BoostingQuery          `json:"boosting"`            // TODO: not implemented
	Match             map[string]interface{}  `json:"match"`               // field: query, field: {query: value, operator: "and|or"}
	MatchBoolPrefix   map[string]interface{}  `json:"match_bool_prefix"`   // field: query, field: {query: value, operator: "and|or"}
	MatchPhrase       map[string]interface{}  `json:"match_phrase"`        // field: query, field: {query: value, operator: "and|or"}
	MatchPhrasePrefix map[string]interface{}  `json:"match_phrase_prefix"` // field: query, field: {query: value, operator: "and|or"}
	MatchAll          interface{}             `json:"match_all"`           // just set or null
	MatchNone         interface{}             `json:"match_none"`          // just set or null
	MultiMatch        *MultiMatchQuery        `json:"multi_match"`         // .
	CombinedFields    *CombinedFieldsQuery    `json:"combined_fields"`     // TODO: not implemented
	QueryString       *QueryStringQuery       `json:"query_string"`        // .
	SimpleQueryString *SimpleQueryStringQuery `json:"simple_query_string"` // .
	Exists            *ExistsQuery            `json:"exists"`              // .
	Ids               *IdsQuery               `json:"ids"`                 // .
	Fuzzy             map[string]interface{}  `json:"fuzzy"`               // field: query, field: {query: value, operator: "and|or"}
	Prefix            map[string]interface{}  `json:"prefix"`              // field: query, field: {query: value, operator: "and|or"}
	Range             *RangeQuery             `json:"range"`               // .
	Term              map[string]interface{}  `json:"term"`                // field: query, field: {query: value, operator: "and|or"}
	Terms             *TermsQuery             `json:"terms"`               // .
	TermsSet          *TermsSetQuery          `json:"terms_set"`           // TODO: not implemented
	Wildcard          map[string]interface{}  `json:"wildcard"`            // field: query, field: {query: value, operator: "and|or"}
	GeoBoundingBox    interface{}             `json:"geo_bounding_box"`    // TODO: not implemented
	GeoDistance       interface{}             `json:"geo_distance"`        // TODO: not implemented
	GeoPolygon        interface{}             `json:"geo_polygon"`         // TODO: not implemented
	GeoShape          interface{}             `json:"geo_shape"`           // TODO: not implemented
}

type BoolQuery struct {
	Should             interface{} `json:"should"`               // query, [query1, query2]
	Must               interface{} `json:"must"`                 // query, [query1, query2]
	MustNot            interface{} `json:"must_not"`             // query, [query1, query2]
	Filter             interface{} `json:"filter"`               // query, [query1, query2]
	MinimumShouldMatch int64       `json:"minimum_should_match"` // only for should
}

type BoostingQuery struct {
	Positive      interface{} `json:"positive"` // singe or multiple queries
	Negative      interface{} `json:"negative"` // singe or multiple queries
	NegativeBoost float64     `json:"negative_boost"`
}

type MatchQuery struct {
	Query          string `json:"query"`
	Analyzer       string `json:"analyzer"`
	Operator       string `json:"operator"`         // or(default), and
	Fuzziness      string `json:"fuzziness"`        // auto(default), 1,2,3,n
	ZeroTermsQuery string `json:"zero_terms_query"` // none(default), all
}

type MatchBoolPrefixQuery struct{}

type MatchPhraseQuery struct{}

type MatchPhrasePrefix struct{}

type MultiMatchQuery struct{}

type CombinedFieldsQuery struct{}

type QueryStringQuery struct{}

type SimpleQueryStringQuery struct{}

type ExistsQuery struct{}

type IdsQuery struct{}

type FuzzyQuery struct{}

type PrefixQuery struct{}

type RangeQuery struct{}

type TermQuery struct{}

type TermsQuery struct{}

type TermsSetQuery struct{}

type WildcardQuery struct{}

type Aggregations struct {
	Avg           AggregationMetric         `json:"avg"`
	Max           AggregationMetric         `json:"max"`
	Min           AggregationMetric         `json:"min"`
	Sum           AggregationMetric         `json:"sum"`
	Count         AggregationMetric         `json:"count"`
	Terms         *AggregationsTerms        `json:"terms"`
	Range         *AggregationRange         `json:"range"`
	DateRange     *AggregationDateRange     `json:"date_range"`
	IPRange       *AggregationIPRange       `json:"ip_range"`       // TODO: not implemented
	Histogram     *AggregationHistogram     `json:"histogram"`      // TODO: not implemented
	DateHistogram *AggregationDateHistogram `json:"date_histogram"` // TODO: not implemented
	Aggregations  map[string]Aggregations   `json:"aggs"`           // nested aggregations
}

type AggregationMetric struct {
	Field string `json:"field"`
}

type AggregationsTerms struct {
	Field string            `json:"field"`
	Size  int64             `json:"size"`
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
	Field    string `json:"field"`
	Interval int64  `json:"interval"`
	Keyed    bool   `json:"keyed"`
}

type AggregationDateHistogram struct {
	Field            string `json:"field"`
	Format           string `json:"format"`            // format key_as_string
	FixedInterval    string `json:"fixed_interval"`    // ms,s,m,h,d
	CalendarInterval string `json:"calendar_interval"` // minute,hour,day,week,month,quarter,year
	Keyed            bool   `json:"keyed"`
}

type Highlight struct {
	NumberOfFragments int64                `json:"number_of_fragments"`
	FragmentSize      int64                `json:"fragment_size"`
	PreTags           []string             `json:"pre_tags"`
	PostTags          []string             `json:"post_tags"`
	Fields            map[string]Highlight `json:"fields"`
}

type Field struct {
	Field  string `json:"field"`
	Format string `json:"format"`
}

type Source struct {
	Enable bool            // enable _source returns, default is true
	Fields map[string]bool // what fields can returns
}

type Sort struct {
	Field string
	Order string // asc, desc
}
