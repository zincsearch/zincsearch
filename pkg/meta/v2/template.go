package v2

type Template struct {
	IndexPatterns []string         `json:"index_patterns"`
	Priority      int              `json:"priority"` // highest priority is chosen
	Template      TemplateTemplate `json:"template"`
}

type TemplateTemplate struct {
	Settings *IndexSettings `json:"settings,omitempty"`
	Mappings *Mappings      `json:"mappings,omitempty"`
}
