package v2

type Mappings struct {
	Properties map[string]Property `json:"properties"`
}

type Property struct {
	Type           string `json:"type"` // text, keyword, date_time, numeric, boolean, geo_point
	Analyzer       string `json:"analyzer"`
	SearchSnalyzer string `json:"search_analyzer"`
	Format         string `json:"format"` // date format yyyy-MM-dd HH:mm:ss || yyyy-MM-dd || epoch_millis
	Index          bool   `json:"index"`
	Store          bool   `json:"store"`
	Sortable       bool   `json:"sortable"`
	Aggregatable   bool   `json:"aggregatable"`
	Highlightable  bool   `json:"highlightable"`
}
