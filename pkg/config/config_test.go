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
		os.Setenv("ZINC_NODE_ID", "x")
		os.Setenv("ZINC_ETCD_ENDPOINTS", "localhost:2379")
	})

	t.Run("check", func(t *testing.T) {
		c := new(config)
		rv := reflect.ValueOf(c).Elem()
		loadConfig(rv)

		assert.Equal(t, "", c.GinMode)
		assert.Equal(t, "4080", c.ServerPort)
		assert.Equal(t, "node", c.ServerMode)
		assert.Equal(t, "x", c.NodeID)
		assert.Equal(t, "./data", c.DataPath)
		assert.Equal(t, true, c.SentryEnable)
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
