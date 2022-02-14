package meta

type Index struct {
	Settings IndexSettings  `json:"settings"`
	Analysis *IndexAnalysis `json:"analysis"`
	Mappings *Mappings      `json:"mappings"`
}

type IndexSettings struct {
	NumberOfShards   int `json:"number_of_shards"`
	NumberOfReplicas int `json:"number_of_replicas"`
}

type IndexAnalysis struct {
}
