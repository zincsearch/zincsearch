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
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
)

// isolate clears every config env var and moves into an empty working
// directory (with an empty ./conf) so only the embedded defaults and
// explicit test inputs apply.
func isolate(t *testing.T) {
	t.Helper()
	clearConfigEnv(t, reflect.TypeFor[config]())
	t.Chdir(t.TempDir())
	assert.NoError(t, os.Mkdir(filepath.Dir(configFile), 0o755))
}

func mustLoad(t *testing.T) *config {
	t.Helper()
	c := new(config)
	if !assert.NoError(t, load(c)) {
		t.FailNow()
	}
	return c
}

func TestConfig(t *testing.T) {
	isolate(t)
	t.Setenv("ZINC_SERVER_MODE", "node")
	t.Setenv("ZINC_NODE_ID", "8")
	t.Setenv("ZINC_ETCD_ENDPOINTS", "localhost:2379")
	t.Setenv("ZINC_MAX_DOCUMENT_SIZE", "1m")
	t.Setenv("ZINC_WAL_SYNC_INTERVAL", "10s")

	t.Run("check", func(t *testing.T) {
		c := mustLoad(t)

		assert.Equal(t, "", c.GinMode)
		assert.Equal(t, "4080", c.ServerPort)
		assert.Equal(t, "node", c.ServerMode)
		assert.Equal(t, 8, c.NodeID)
		assert.Equal(t, "./data", c.DataPath)
		assert.Equal(t, false, c.SentryEnable)
		assert.Equal(t, "", c.SentryDSN)
		assert.Equal(t, false, c.TelemetryEnable)
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
			{value: "2048576", expect: 2048576},
			{value: "1k", expect: 1000},
			{value: "1kb", expect: 1000},
			{value: "1m", expect: 1000000},
			{value: "1mb", expect: 1000000},
			{value: "1g", expect: 1000000000},
			{value: "1gb", expect: 1000000000},
			{value: "1G", expect: 1000000000},
			{value: "1GB", expect: 1000000000},
		}
		for _, v := range tests {
			t.Setenv("ZINC_MAX_DOCUMENT_SIZE", v.value)
			assert.Equal(t, v.expect, mustLoad(t).MaxDocumentSize)
		}

		dt := []struct {
			value  string
			expect time.Duration
		}{
			{value: "1", expect: time.Nanosecond},
			{value: "1ns", expect: time.Nanosecond},
			{value: "1s", expect: time.Second},
			{value: "1m", expect: time.Minute},
		}
		for _, v := range dt {
			t.Setenv("ZINC_WAL_SYNC_INTERVAL", v.value)
			assert.Equal(t, v.expect, mustLoad(t).WalSyncInterval)
		}
	})
}

func TestSentryDSNOverride(t *testing.T) {
	isolate(t)
	customDSN := "https://secretToken.my.sentry.com/1234"
	t.Setenv("ZINC_SENTRY_DSN", customDSN)
	assert.Equal(t, customDSN, mustLoad(t).SentryDSN)
}

func TestTOMLDefaults(t *testing.T) {
	c := new(config)
	decoder := toml.NewDecoder(bytes.NewReader(defaultConfig)).DisallowUnknownFields()
	if !assert.NoError(t, decoder.Decode(c)) {
		return
	}

	assert.Equal(t, "4080", c.ServerPort)
	assert.Equal(t, "node", c.ServerMode)
	assert.Equal(t, 1, c.NodeID)
	assert.Equal(t, "./data", c.DataPath)
	assert.Equal(t, "bolt", c.MetadataStorage)
	assert.Equal(t, "zstd", c.IceCompressor)
	assert.False(t, c.SentryEnable)
	assert.False(t, c.TelemetryEnable)
	assert.False(t, c.ProfilerEnable)
	assert.False(t, c.PrometheusEnable)
	assert.False(t, c.EnableTextKeywordMapping)
	assert.True(t, c.EnableSQLMappings)
	assert.Equal(t, 1024, c.BatchSize)
	assert.Equal(t, 10000, c.MaxResults)
	assert.Equal(t, 1000, c.AggregationTermsSize)
	assert.Equal(t, 1000000, c.MaxDocumentSize)
	assert.Equal(t, time.Second, c.WalSyncInterval)
	assert.False(t, c.WalRedoLogNoSync)
	assert.True(t, c.ZincSwaggerEnable)
	assert.Equal(t, "debug", c.LogLevel)
	assert.Equal(t, cluster{Name: "ZincCluster"}, c.Cluster)
	assert.Equal(t, shard{GoroutineNum: 3, Num: 3, MaxSize: 1073741824}, c.Shard)
	assert.Equal(t, etcd{Prefix: "/zinc"}, c.Etcd)
	assert.Equal(t, gse{
		EnableStop: true,
		EnableHMM:  true,
		DictEmbed:  "small",
		DictPath:   "./plugins/gse/dict",
	}, c.Plugin.GSE)
}

func TestTOMLEnvironmentOverrides(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
	}{
		{name: "unset", env: nil},
		{name: "empty", env: map[string]string{
			"ZINC_SERVER_PORT": "", "ZINC_WAL_SYNC_INTERVAL": "",
			"ZINC_PLUGIN_GSE_ENABLE_STOP": "",
		}},
		{name: "overridden", env: map[string]string{
			"ZINC_SERVER_PORT": "9090", "ZINC_NODE_ID": "0",
			"ZINC_SENTRY": "false", "ZINC_WAL_SYNC_INTERVAL": "25ms",
			"ZINC_MAX_DOCUMENT_SIZE": "2mb", "ZINC_SHARD_MAX_SIZE": "2048",
			"ZINC_ETCD_ENDPOINTS":         "one:2379,two:2379",
			"ZINC_PLUGIN_GSE_ENABLE_STOP": "false",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolate(t)
			for key, value := range tt.env {
				t.Setenv(key, value)
			}
			want := new(config)
			if !assert.NoError(t, toml.Unmarshal(defaultConfig, want)) {
				return
			}
			if tt.name == "overridden" {
				want.ServerPort = "9090"
				want.NodeID = 0
				want.SentryEnable = false
				want.WalSyncInterval = 25 * time.Millisecond
				want.MaxDocumentSize = 2000000
				want.Shard.MaxSize = 2048
				want.Etcd.Endpoints = []string{"one:2379", "two:2379"}
				want.Plugin.GSE.EnableStop = false
			}
			assert.Equal(t, want, mustLoad(t))
		})
	}
}

func clearConfigEnv(t *testing.T, rt reflect.Type) {
	t.Helper()
	for field := range rt.Fields() {
		if field.Type.Kind() == reflect.Struct {
			clearConfigEnv(t, field.Type)
		} else if tag := field.Tag.Get("toml"); tag != "" {
			assert.Equal(t, strings.ToLower(tag), tag, "toml tag must be lowercase")
			key := strings.ToUpper(tag)
			t.Setenv(key, "")
			assert.NoError(t, os.Unsetenv(key))
		}
	}
}

func TestLoadConfTOML(t *testing.T) {
	isolate(t)
	const contents = `zinc_first_admin_user = "toml-admin"
zinc_first_admin_password = "test-only-password"
zinc_server_port = "9080"
zinc_enable_sql_mappings = false
[Shard]
zinc_shard_num = 5
`
	if !assert.NoError(t, os.WriteFile(configFile, []byte(contents), 0o600)) {
		return
	}
	c := mustLoad(t)
	assert.Equal(t, "toml-admin", c.FirstAdminUser)
	assert.Equal(t, "test-only-password", c.FirstAdminPassword)
	assert.Equal(t, "9080", c.ServerPort)
	assert.False(t, c.EnableSQLMappings)
	assert.Equal(t, int64(5), c.Shard.Num)
	assert.Equal(t, 3, c.Shard.GoroutineNum)
	assert.Equal(t, time.Second, c.WalSyncInterval)

	// process environment wins over the TOML file
	t.Setenv("ZINC_FIRST_ADMIN_PASSWORD", "env-test-password")
	t.Setenv("ZINC_ENABLE_SQL_MAPPINGS", "true")
	c = mustLoad(t)
	assert.Equal(t, "env-test-password", c.FirstAdminPassword)
	assert.True(t, c.EnableSQLMappings)
	assert.Equal(t, "9080", c.ServerPort)
}

func TestLoadIgnoresOtherSources(t *testing.T) {
	isolate(t)
	// Only ./conf/zinc.toml is read; ZINC_CONFIG_FILE, .env, ./zinc.toml and default.toml are not config sources.
	assert.Equal(t, filepath.Join("conf", "zinc.toml"), filepath.Clean(configFile))
	for _, name := range []string{"other.toml", "zinc.toml", "default.toml", "conf/default.toml"} {
		assert.NoError(t, os.WriteFile(name, []byte(`zinc_server_port = "1"`+"\n"), 0o600))
	}
	assert.NoError(t, os.WriteFile(".env", []byte("ZINC_SERVER_PORT=2\n"), 0o600))
	t.Setenv("ZINC_CONFIG_FILE", "other.toml")
	assert.Equal(t, "4080", mustLoad(t).ServerPort)
}

func TestLoadInvalidTOML(t *testing.T) {
	isolate(t)
	// Go field names are not accepted as keys; only zinc_* names are (matched case-insensitively by go-toml).
	for _, content := range []string{`zinc_server_port = [`, `UnknownSetting = true`, `ServerPort = "1"`, `zinc_first_admin_password = 123`} {
		assert.NoError(t, os.WriteFile(configFile, []byte(content), 0o600))
		err := load(new(config))
		assert.ErrorContains(t, err, "invalid TOML configuration")
		assert.NotContains(t, err.Error(), content)
	}
	assert.NoError(t, os.Mkdir("dir", 0o755))
	t.Chdir("dir")
	assert.NoError(t, os.MkdirAll(configFile, 0o755))
	assert.ErrorContains(t, load(new(config)), "read config")
}
