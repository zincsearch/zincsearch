package analyzer

// CharFilter
type CharFilter struct {
	Type string `json:"type"`
}

// ascii_folding
type AsciiFoldingCharFilter struct{}

// html
type HtmlCharFilter struct{}

// regexp
type RegexpCharFilter struct {
	CharFilter
	Pattern     string `json:"pattern"`
	Replacement string `json:"replacement"`
	Flags       string `json:"flags,omitempty"`
}

// mapping
type MappingCharFilter struct {
	CharFilter
	Mappings []string `json:"mappings"`
}

// zero_width_non_joiner
type ZeroWidthNonJoinerCharFilter struct{}
