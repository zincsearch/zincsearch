package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"
	"github.com/rs/zerolog/log"

	meta "github.com/zinclabs/zinc/pkg/meta/v2"
	parser "github.com/zinclabs/zinc/pkg/uquery/v2"
)

func MultiSearchV2(indexName string, query *meta.ZincQuery) (*meta.SearchResponse, error) {
	var mappings *meta.Mappings
	var analyzers map[string]*analysis.Analyzer
	var readers []*bluge.Reader
	for name, index := range ZINC_INDEX_LIST {
		if indexName == "" || (indexName != "" && strings.HasPrefix(name, indexName[:len(indexName)-1])) {
			reader, _ := index.Writer.Reader()
			readers = append(readers, reader)
			if mappings == nil {
				mappings = index.CachedMappings
				analyzers = index.CachedAnalyzers
			}
		}
	}

	if len(readers) == 0 {
		return nil, fmt.Errorf("core.MultiSearchV2: error accessing reader: no index found")
	}

	searchRequest, err := parser.ParseQueryDSL(query, mappings, analyzers)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var cancel context.CancelFunc
	if query.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(query.Timeout)*time.Second)
		defer cancel()
	}

	dmi, err := bluge.MultiSearch(ctx, searchRequest, readers...)
	if err != nil {
		log.Printf("core.MultiSearchV2: error executing search: %s", err.Error())
		if err == context.DeadlineExceeded {
			return &meta.SearchResponse{
				TimedOut: true,
				Error:    err.Error(),
				Hits:     meta.Hits{Hits: []meta.Hit{}},
			}, nil
		}
		return nil, err
	}

	return searchV2(dmi, query, mappings)
}
