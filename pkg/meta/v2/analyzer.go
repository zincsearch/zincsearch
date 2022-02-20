package v2

type Analyzer struct {
	CharFilter  *CharFilter  `json:"char_filter"`
	Tokenizer   *Tokenizer   `json:"tokenizer"`
	TokenFilter *TokenFilter `json:"token_filter"`
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
