package v2

import "github.com/prabhatsharma/zinc/pkg/meta/v2/analyzer"

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
	Analyzer    map[string]*analyzer.Analyzer    `json:"analyzer,omitempty"`
	CharFilter  map[string]interface{}           `json:"char_filter,omitempty"`
	Tokenizer   map[string]*analyzer.Tokenizer   `json:"tokenizer,omitempty"`
	TokenFilter map[string]*analyzer.TokenFilter `json:"token_filter,omitempty"`
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
