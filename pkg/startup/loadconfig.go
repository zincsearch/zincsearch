package startup

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const (
	DEFAULT_BATCH_SIZE             = 1000
	DEFAULT_MAX_RESULTS            = 10000
	DEFAULT_AGGREGATION_TERMS_SIZE = 100
)

var batchSize = DEFAULT_BATCH_SIZE
var maxResults = DEFAULT_MAX_RESULTS
var aggregationTermsSize = DEFAULT_AGGREGATION_TERMS_SIZE

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Info().Msg("Error loading .env file")
	}

	var vs string
	var vi int
	vs = os.Getenv("ZINC_BATCH_SIZE")
	if vs != "" {
		if vi, err = strconv.Atoi(vs); err == nil {
			batchSize = vi
		}
	}

	vs = os.Getenv("ZINC_MAX_RESULTS")
	if vs != "" {
		if vi, err = strconv.Atoi(vs); err == nil {
			maxResults = vi
		}
	}

	vs = os.Getenv("ZINC_AGGREGATION_TERMS_SIZE")
	if vs != "" {
		if vi, err = strconv.Atoi(vs); err == nil {
			aggregationTermsSize = vi
		}
	}

}

func LoadBatchSize() int {
	return batchSize
}

func LoadMaxResults() int {
	return maxResults
}

func LoadAggregationTermsSize() int {
	return aggregationTermsSize
}
