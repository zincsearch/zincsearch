package v2

type Analyzer struct {
	CharFilter  []string `json:"char_filter,omitempty"`
	Tokenizer   string   `json:"tokenizer,omitempty"`
	TokenFilter []string `json:"token_filter,omitempty"`
	Filter      []string `json:"filter,omitempty"` // compatibility with es, alias for TokenFilter

	// options for compatible
	Type      string   `json:"type,omitempty"`
	Pattern   string   `json:"pattern,omitempty"`   // for type=pattern
	Lowercase bool     `json:"lowercase,omitempty"` // for type=pattern
	Stopwords []string `json:"stopwords,omitempty"` // for type=pattern,standard,stop
}

type Tokenizer struct {
	Type string `json:"type"`
}
type TokenFilter struct {
	Type string `json:"type"`
}
