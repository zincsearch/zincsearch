package v2

type ZincQuery struct {
	Query          interface{}             `json:"query"`
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
