package v2

type Analyzer struct {
	CharFilter  []string `json:"char_filter,omitempty"`
	Tokenizer   string   `json:"tokenizer,omitempty"`
	TokenFilter []string `json:"token_filter,omitempty"`
}

type CharFilter struct {
	Type string `json:"type"`
}

type Tokenizer struct {
	Type string `json:"type"`
}
type TokenFilter struct {
	Type string `json:"type"`
}
