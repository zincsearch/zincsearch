package v2

type Index struct {
	Name        string         `json:"name"`
	DocsCount   int64          `json:"docs_count"`
	StorageSize int64          `json:"storage_size"`
	StorageType string         `json:"storage_type"`
	Settings    *IndexSettings `json:"settings,omitempty"`
	Mappings    *Mappings      `json:"mappings,omitempty"`
}

type IndexSettings struct {
	NumberOfShards   int            `json:"number_of_shards,omitempty"`
	NumberOfReplicas int            `json:"number_of_replicas,omitempty"`
	Analysis         *IndexAnalysis `json:"analysis,omitempty"`
}

type IndexAnalysis struct {
	Analyzer    map[string]*Analyzer   `json:"analyzer,omitempty"`
	CharFilter  map[string]interface{} `json:"char_filter,omitempty"`
	Tokenizer   map[string]interface{} `json:"tokenizer,omitempty"`
	TokenFilter map[string]interface{} `json:"token_filter,omitempty"`
	Filter      map[string]interface{} `json:"filter,omitempty"` // compatibility with es, alias for TokenFilter
}

type SortIndex []*Index

func (t SortIndex) Len() int {
	return len(t)
}

func (t SortIndex) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t SortIndex) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}
