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

package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("prepare", func(t *testing.T) {
		os.Setenv("ZINC_SERVER_MODE", "node")
		os.Setenv("ZINC_NODE_ID", "8")
		os.Setenv("ZINC_ETCD_ENDPOINTS", "localhost:2379")
	})

	t.Run("check", func(t *testing.T) {
		c := new(config)
		rv := reflect.ValueOf(c).Elem()
		loadConfig(rv)

		assert.Equal(t, "", c.GinMode)
		assert.Equal(t, "4080", c.ServerPort)
		assert.Equal(t, "node", c.ServerMode)
		assert.Equal(t, 8, c.NodeID)
		assert.Equal(t, "./data", c.DataPath)
		assert.Equal(t, true, c.SentryEnable)
		assert.Equal(t, "https://15b6d9b8be824b44896f32b0234c32b7@o1218932.ingest.sentry.io/6360942", c.SentryDSN) // Add check for default value
		assert.Equal(t, true, c.TelemetryEnable)
		assert.Equal(t, false, c.PrometheusEnable)

		assert.Equal(t, 1024, c.BatchSize)
		assert.Equal(t, 10000, c.MaxResults)
		assert.Equal(t, 1000, c.AggregationTermsSize)

		assert.Equal(t, []string{"localhost:2379"}, c.Etcd.Endpoints)

		assert.Equal(t, "", c.S3.Bucket)
		assert.Equal(t, "", c.MinIO.Endpoint)

		assert.Equal(t, false, c.Plugin.GSE.Enable)
		assert.Equal(t, "small", c.Plugin.GSE.DictEmbed)
		assert.Equal(t, "./plugins/gse/dict", c.Plugin.GSE.DictPath)
	})
}

func TestSentryDSNOverride(t *testing.T) {
	customDSN := "https://secretToken.my.sentry.com/1234"

	t.Run("prepare", func(t *testing.T) {
		os.Setenv("ZINC_SENTRY_DSN", customDSN)
	})

	t.Run("check", func(t *testing.T) {
		c := new(config)
		rv := reflect.ValueOf(c).Elem()
		loadConfig(rv)

		assert.Equal(t, customDSN, c.SentryDSN)
	})
}

func TestS3Override(t *testing.T) {
	bucket := "zinc-dev-misc1"

	t.Run("prepare", func(t *testing.T) {
		os.Setenv("ZINC_S3_BUCKET", bucket)
	})

	t.Run("check", func(t *testing.T) {
		c := new(config)
		rv := reflect.ValueOf(c).Elem()
		loadConfig(rv)

		assert.Equal(t, bucket, c.S3.Bucket)
	})
}
