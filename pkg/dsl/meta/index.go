package meta

type Index struct {
	Settings IndexSettings  `json:"settings"`
	Analysis *IndexAnalysis `json:"analysis,omitempty"`
	Mappings *Mappings      `json:"mappings,omitempty"`
}

type IndexSettings struct {
	NumberOfShards   int `json:"number_of_shards"`
	NumberOfReplicas int `json:"number_of_replicas"`
}

type IndexAnalysis struct {
}

func NewIndex() *Index {
	return &Index{
		Settings: IndexSettings{
			NumberOfReplicas: 1,
			NumberOfShards:   3,
		},
	}
}
