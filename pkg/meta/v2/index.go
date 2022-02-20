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
	Analyzer    map[string]*Analyzer    `json:"analyzer"`
	CharFilter  map[string]*CharFilter  `json:"char_filter"`
	Tokenizer   map[string]*Tokenizer   `json:"tokenizer"`
	TokenFilter map[string]*TokenFilter `json:"token_filter"`
}

func NewIndex() *Index {
	return &Index{
		Settings: &IndexSettings{
			NumberOfReplicas: 1,
			NumberOfShards:   3,
		},
	}
}
