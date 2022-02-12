package startup

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	DEFAULT_BATCH_SIZE  = 1000
	DEFAULT_MAX_RESULTS = 10000
)

var batchSize = DEFAULT_BATCH_SIZE
var maxResults = DEFAULT_MAX_RESULTS

func init() {
	godotenv.Load()

	var vs string
	var vi int
	var err error
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

}

func LoadBatchSize() int {
	return batchSize
}

func LoadMaxResults() int {
	return maxResults
}
