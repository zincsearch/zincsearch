package v2

type Mappings struct {
	Properties map[string]Property `json:"properties,omitempty"`
}

type Property struct {
	Type           string `json:"type"` // text, keyword, time, numeric, boolean, geo_point
	Analyzer       string `json:"analyzer,omitempty"`
	SearchAnalyzer string `json:"search_analyzer,omitempty"`
	Format         string `json:"format,omitempty"` // date format yyyy-MM-dd HH:mm:ss || yyyy-MM-dd || epoch_millis
	Index          bool   `json:"index"`
	Store          bool   `json:"store"`
	Sortable       bool   `json:"sortable"`
	Aggregatable   bool   `json:"aggregatable"`
	Highlightable  bool   `json:"highlightable"`
}

func NewMappings() *Mappings {
	return &Mappings{
		Properties: make(map[string]Property),
	}
}

func NewProperty(typ string) Property {
	p := Property{
		Type:           typ,
		Analyzer:       "",
		SearchAnalyzer: "",
		Format:         "",
		Index:          true,
		Store:          false,
		Sortable:       true,
		Aggregatable:   true,
		Highlightable:  false,
	}
	if typ == "text" {
		p.Sortable = false
		p.Aggregatable = false
	}

	return p
}
