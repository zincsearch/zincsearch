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
	_ "embed"
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vcaesar/ice/compress"
)

type config struct {
	FirstAdminUser            string        `toml:"zinc_first_admin_user"`
	FirstAdminPassword        string        `toml:"zinc_first_admin_password"`
	GinMode                   string        `toml:"gin_mode"`
	ServerPort                string        `toml:"zinc_server_port"`
	ServerAddress             string        `toml:"zinc_server_address"`
	ServerTLSCertificateFile  string        `toml:"zinc_server_tls_certificate_file"`
	ServerTLSKeyFile          string        `toml:"zinc_server_tls_key_file"`
	ServerMode                string        `toml:"zinc_server_mode"`
	NodeID                    int           `toml:"zinc_node_id"`
	DataPath                  string        `toml:"zinc_data_path"`
	MetadataStorage           string        `toml:"zinc_metadata_storage"`
	IceCompressor             string        `toml:"zinc_ice_compressor"`
	SentryEnable              bool          `toml:"zinc_sentry"`
	SentryDSN                 string        `toml:"zinc_sentry_dsn"`
	ProfilerEnable            bool          `toml:"zinc_profiler"`
	ProfilerServer            string        `toml:"zinc_profiler_server"`
	ProfilerAPIKey            string        `toml:"zinc_profiler_api_key"`
	ProfilerFriendlyProfileID string        `toml:"zinc_profiler_friendly_profile_id"`
	TelemetryEnable           bool          `toml:"zinc_telemetry"`
	PrometheusEnable          bool          `toml:"zinc_prometheus_enable"`
	EnableTextKeywordMapping  bool          `toml:"zinc_enable_text_keyword_mapping"`
	EnableSQLMappings         bool          `toml:"zinc_enable_sql_mappings"`
	BatchSize                 int           `toml:"zinc_batch_size"`
	MaxResults                int           `toml:"zinc_max_results"`
	AggregationTermsSize      int           `toml:"zinc_aggregation_terms_size"`
	MaxDocumentSize           int           `toml:"zinc_max_document_size"`   // Max size for a single document . Default = 1 MB = 1024 * 1024
	WalSyncInterval           time.Duration `toml:"zinc_wal_sync_interval"`   // sync wal to disk, 1s, 10ms
	WalRedoLogNoSync          bool          `toml:"zinc_wal_redolog_no_sync"` // control sync after every write
	ZincSwaggerEnable         bool          `toml:"zinc_swagger_enable"`
	LogLevel                  string        `toml:"zinc_log_level"`
	Cluster                   cluster
	Shard                     shard
	Etcd                      etcd
	Plugin                    plugin
}

type cluster struct {
	Name string `toml:"zinc_cluster_name"`
}

type shard struct {
	// control goroutine number for read
	GoroutineNum int `toml:"zinc_shard_goroutine_num"`
	// DefaultNum is the default number of shards.
	Num int64 `toml:"zinc_shard_num"`
	// MaxSize is the maximum size limit for one shard, or will create a new shard.
	MaxSize uint64 `toml:"zinc_shard_max_size"`
}

type etcd struct {
	Endpoints []string `toml:"zinc_etcd_endpoints"`
	Prefix    string   `toml:"zinc_etcd_prefix"`
	Username  string   `toml:"zinc_etcd_username"`
	Password  string   `toml:"zinc_etcd_password"`
}

type plugin struct {
	ES  elasticsearch
	GSE gse
}

type elasticsearch struct {
	Version string `toml:"zinc_plugin_es_version"`
}

type gse struct {
	Enable     bool   `toml:"zinc_plugin_gse_enable"`
	EnableStop bool   `toml:"zinc_plugin_gse_enable_stop"`
	EnableHMM  bool   `toml:"zinc_plugin_gse_enable_hmm"`
	DictEmbed  string `toml:"zinc_plugin_gse_dict_embed"`
	DictPath   string `toml:"zinc_plugin_gse_dict_path"`
}

// configFile is the only config file read at startup, resolved against the working directory.
const configFile = "conf/zinc.toml"

//go:embed default.toml
var defaultConfig []byte

var Global = new(config)

func init() {
	if err := load(Global); err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	// set the log level
	logLevel, logLevelErr := zerolog.ParseLevel(Global.LogLevel)
	if logLevelErr != nil {
		log.Error().Err(logLevelErr).Msg("ZINC_LOG_LEVEL is not valid")
	}
	zerolog.SetGlobalLevel(logLevel)

	// configure gin
	if Global.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// check data path
	testPath := path.Join(Global.DataPath, "_test_")
	if err := os.MkdirAll(testPath, 0o755); err != nil {
		log.Fatal().Err(err).Msg("ZINC_DATA_PATH is not writable")
	}
	if err := os.Remove(testPath); err != nil {
		log.Fatal().Err(err).Msg("ZINC_DATA_PATH is not writable")
	}

	// configure ice compress algorithm
	switch strings.ToUpper(Global.IceCompressor) {
	case "SNAPPY":
		compress.Algorithm = compress.SNAPPY
	case "S2":
		compress.Algorithm = compress.S2
	case "ZSTD":
		compress.Algorithm = compress.ZSTD
	}
}

// load fills c from the embedded defaults, then ./conf/zinc.toml when present,
// then non-empty process environment variables. TOML keys are lowercase and
// map to the same names uppercased in the environment.
func load(c *config) error {
	if err := toml.Unmarshal(defaultConfig, c); err != nil {
		return fmt.Errorf("decode built-in defaults: %w", err)
	}
	data, err := os.ReadFile(configFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read config %s: %w", configFile, err)
	}
	if err == nil {
		if err := toml.NewDecoder(bytes.NewReader(data)).DisallowUnknownFields().Decode(c); err != nil {
			// Decoder errors can contain configuration values, including passwords.
			return fmt.Errorf("invalid TOML configuration in %s", configFile)
		}
	}
	loadEnv(reflect.ValueOf(c).Elem())
	return nil
}

func loadEnv(rv reflect.Value) {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)
		if ft.Type.Kind() == reflect.Struct {
			loadEnv(fv)
			continue
		}
		if tag := ft.Tag.Get("toml"); tag != "" {
			setField(fv, strings.ToUpper(tag))
		}
	}
}

func setField(field reflect.Value, tag string) {
	if tag == "" {
		return
	}
	v := os.Getenv(tag)
	if v == "" {
		return
	}
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, ok := field.Interface().(time.Duration)
		var (
			vi  int64
			err error
		)
		switch ok {
		case true:
			d, e := time.ParseDuration(v)
			if e != nil && strings.Contains(e.Error(), "time: missing unit in duration") {
				vi, err = strconv.ParseInt(v, 10, 64)
			} else {
				vi, err = int64(d), e
			}

		default:
			vi, err = units.FromHumanSize(v)
		}
		if err != nil {
			log.Fatal().Err(err).Msgf("env %s is not int", tag)
		}

		field.SetInt(int64(vi))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		vi, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			log.Fatal().Err(err).Msgf("env %s is not uint", tag)
		}
		field.SetUint(uint64(vi))
	case reflect.Bool:
		vi, err := strconv.ParseBool(v)
		if err != nil {
			log.Fatal().Err(err).Msgf("env %s is not bool", tag)
		}
		field.SetBool(vi)
	case reflect.String:
		field.SetString(v)
	case reflect.Slice:
		vs := strings.Split(v, ",")
		field.Set(reflect.ValueOf(vs))
		field.SetLen(len(vs))
	}
}
