package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert.Equal(t, "1", Config.NodeID)
	assert.Equal(t, "./data", Config.DataPath)
	assert.Equal(t, "", Config.GinMode)
	assert.Equal(t, true, Config.SentryEnable)
	assert.Equal(t, true, Config.TelemetryEnable)
	assert.Equal(t, false, Config.PrometheusEnable)

	assert.Equal(t, 1024, Config.BatchSize)
	assert.Equal(t, 10000, Config.MaxResults)
	assert.Equal(t, 1000, Config.AggregationTermsSize)

	assert.Equal(t, "", Config.S3.Bucket)
	assert.Equal(t, "", Config.MinIO.Endpoint)

	assert.Equal(t, false, Config.Plugin.GSE.Enable)
	assert.Equal(t, "small", Config.Plugin.GSE.DictEmbed)
	assert.Equal(t, "./plugins/gse/dict", Config.Plugin.GSE.DictPath)
}
