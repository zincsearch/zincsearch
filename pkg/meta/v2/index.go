package v2

type Index struct {
	Settings *IndexSettings `json:"settings,omitempty"`
	Mappings *Mappings      `json:"mappings,omitempty"`
}

type IndexSettings struct {
	NumberOfShards   int            `json:"number_of_shards"`
	NumberOfReplicas int            `json:"number_of_replicas"`
	Analysis         *IndexAnalysis `json:"analysis,omitempty"`
}

type IndexAnalysis struct {
	Analyzer    map[string]*Analyzer   `json:"analyzer,omitempty"`
	CharFilter  map[string]interface{} `json:"char_filter,omitempty"`
	Tokenizer   map[string]interface{} `json:"tokenizer,omitempty"`
	TokenFilter map[string]interface{} `json:"token_filter,omitempty"`
	Filter      map[string]interface{} `json:"filter,omitempty"` // compatibility with es, alias for TokenFilter
}

func NewIndex() *Index {
	return &Index{
		Settings: NewIndexSettings(),
	}
}

func NewIndexSettings() *IndexSettings {
	return &IndexSettings{
		NumberOfShards:   3,
		NumberOfReplicas: 1,
	}
}
