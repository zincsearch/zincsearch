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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert.Equal(t, "", Global.GinMode)
	assert.Equal(t, "4080", Global.ServerPort)
	// assert.Equal(t, "node", Global.ServerMode)
	assert.Equal(t, "1", Global.NodeID)
	assert.Equal(t, "./data", Global.DataPath)
	assert.Equal(t, true, Global.SentryEnable)
	assert.Equal(t, true, Global.TelemetryEnable)
	assert.Equal(t, false, Global.PrometheusEnable)

	assert.Equal(t, 1024, Global.BatchSize)
	assert.Equal(t, 10000, Global.MaxResults)
	assert.Equal(t, 1000, Global.AggregationTermsSize)

	// assert.Equal(t, []string(nil), Global.Etcd.Endpoints)

	assert.Equal(t, "", Global.S3.Bucket)
	assert.Equal(t, "", Global.MinIO.Endpoint)

	assert.Equal(t, false, Global.Plugin.GSE.Enable)
	assert.Equal(t, "small", Global.Plugin.GSE.DictEmbed)
	assert.Equal(t, "./plugins/gse/dict", Global.Plugin.GSE.DictPath)
}
