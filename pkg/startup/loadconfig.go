/* Copyright 2022 Zinc Labs Inc. and Contributors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

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
	DEFAULT_AGGREGATION_TERMS_SIZE = 1000
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
