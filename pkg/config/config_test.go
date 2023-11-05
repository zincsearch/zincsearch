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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("prepare", func(t *testing.T) {
		os.Setenv("ZINC_SERVER_MODE", "node")
		os.Setenv("ZINC_NODE_ID", "8")
		os.Setenv("ZINC_ETCD_ENDPOINTS", "localhost:2379")
		os.Setenv("ZINC_MAX_DOCUMENT_SIZE", "1m")
		os.Setenv("ZINC_WAL_SYNC_INTERVAL", "10s")
	})

	t.Run("check", func(t *testing.T) {
		c := new(config)
		loadConfig(reflect.ValueOf(c).Elem())

		assert.Equal(t, "", c.GinMode)
		assert.Equal(t, "4080", c.ServerPort)
		assert.Equal(t, "node", c.ServerMode)
		assert.Equal(t, 8, c.NodeID)
		assert.Equal(t, "./data", c.DataPath)
		assert.Equal(t, true, c.SentryEnable)
		assert.Equal(t, "https://15b6d9b8be824b44896f32b0234c32b7@o1218932.ingest.sentry.io/6360942", c.SentryDSN) // Add check for default value
		assert.Equal(t, true, c.TelemetryEnable)
		assert.Equal(t, false, c.PrometheusEnable)
		assert.Equal(t, 1000000, c.MaxDocumentSize)

		assert.Equal(t, 1024, c.BatchSize)
		assert.Equal(t, 10000, c.MaxResults)
		assert.Equal(t, 1000, c.AggregationTermsSize)

		assert.Equal(t, 10*time.Second, c.WalSyncInterval)

		assert.Equal(t, []string{"localhost:2379"}, c.Etcd.Endpoints)

		assert.Equal(t, false, c.Plugin.GSE.Enable)
		assert.Equal(t, "small", c.Plugin.GSE.DictEmbed)
		assert.Equal(t, "./plugins/gse/dict", c.Plugin.GSE.DictPath)
	})

	t.Run("human check", func(t *testing.T) {
		tests := []struct {
			value  string
			expect int
		}{
			{
				value:  "2048576",
				expect: 2048576,
			},
			{
				value:  "1k",
				expect: 1000,
			},
			{
				value:  "1kb",
				expect: 1000,
			},
			{
				value:  "1m",
				expect: 1000000,
			},
			{
				value:  "1mb",
				expect: 1000000,
			},
			{
				value:  "1g",
				expect: 1000000000,
			},
			{
				value:  "1gb",
				expect: 1000000000,
			},
			{
				value:  "1G",
				expect: 1000000000,
			},
			{
				value:  "1GB",
				expect: 1000000000,
			},
		}
		for _, v := range tests {
			os.Setenv("ZINC_MAX_DOCUMENT_SIZE", v.value)

			c := new(config)
			loadConfig(reflect.ValueOf(c).Elem())
			assert.Equal(t, c.MaxDocumentSize, v.expect)
		}

		dt := []struct {
			value  string
			expect time.Duration
		}{
			{
				value:  "1",
				expect: time.Nanosecond,
			},
			{
				value:  "1ns",
				expect: time.Nanosecond,
			},
			{
				value:  "1s",
				expect: time.Second,
			},
			{
				value:  "1m",
				expect: time.Minute,
			},
		}
		for _, v := range dt {
			os.Setenv("ZINC_WAL_SYNC_INTERVAL", v.value)

			c := new(config)
			loadConfig(reflect.ValueOf(c).Elem())
			assert.Equal(t, c.WalSyncInterval, v.expect)
		}
	})
}

func TestSentryDSNOverride(t *testing.T) {
	customDSN := "https://secretToken.my.sentry.com/1234"

	t.Run("prepare", func(t *testing.T) {
		os.Setenv("ZINC_SENTRY_DSN", customDSN)
	})

	t.Run("check", func(t *testing.T) {
		c := new(config)
		loadConfig(reflect.ValueOf(c).Elem())

		assert.Equal(t, customDSN, c.SentryDSN)
	})
}
