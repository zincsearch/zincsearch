package core

import (
	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/rs/zerolog/log"

	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	"github.com/prabhatsharma/zinc/pkg/plugin"
)

const (
	_ = iota
	NotCompatibleNewIndexMeta
	UseNewIndexMeta
)

var ZINC_INDEX_LIST map[string]*Index
var ZINC_SYSTEM_INDEX_LIST map[string]*Index

func init() {
	// need load plugin before load index
	plugin.Load()

	var err error
	ZINC_SYSTEM_INDEX_LIST, err = LoadZincSystemIndexes()
	if err != nil {
		log.Fatal().Msgf("Error loading system index: %s", err.Error())
	}

	ZINC_INDEX_LIST, _ = LoadZincIndexesFromMeta()
	if err != nil {
		log.Error().Msgf("Error loading user index: %s", err.Error())
	}

	// DEPRECATED compatibility with old code < v0.1.7
	if len(ZINC_INDEX_LIST) == 0 {
		log.Error().Bool("deprecated", true).Msg("Loading user indexes for old version...")
		ZINC_INDEX_LIST, _ = LoadZincIndexesFromDisk()
		s3List, _ := LoadZincIndexesFromS3()
		for k, v := range s3List {
			ZINC_INDEX_LIST[k] = v
		}
		minioList, _ := LoadZincIndexesFromMinIO()
		for k, v := range minioList {
			ZINC_INDEX_LIST[k] = v
		}
		// store index for new version
		for _, index := range ZINC_INDEX_LIST {
			StoreIndex(index)
		}
	}
}

type Index struct {
	Name            string                        `json:"name"`
	IndexType       string                        `json:"index_type"`   // "system" or "user"
	StorageType     string                        `json:"storage_type"` // disk, memory, s3
	Size            float64                       `json:"size"`         // cached size of the index
	Mappings        map[string]interface{}        `json:"mappings"`
	Settings        *meta.IndexSettings           `json:"settings"`
	CachedAnalyzers map[string]*analysis.Analyzer `json:"-"`
	CachedMappings  *meta.Mappings                `json:"-"`
	Writer          *bluge.Writer                 `json:"-"`
}

type IndexTemplate struct {
	Name          string         `json:"name"`
	IndexTemplate *meta.Template `json:"index_template"`
}
