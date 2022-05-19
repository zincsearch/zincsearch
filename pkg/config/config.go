package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type config struct {
	NodeID           string `env:"ZINC_NODE_ID,default=1"`
	DataPath         string `env:"ZINC_DATA_PATH,default=./data"`
	GinMode          string `env:"GIN_MODE"`
	SentryEnable     bool   `env:"ZINC_SENTRY,default=true"`
	TelemetryEnable  bool   `env:"ZINC_TELEMETRY,default=true"`
	PrometheusEnable bool   `env:"ZINC_PROMETHEUS_ENABLE,default=false"`

	BatchSize            int `env:"ZINC_BATCH_SIZE,default=1024"`
	MaxResults           int `env:"ZINC_MAX_RESULTS,default=10000"`
	AggregationTermsSize int `env:"ZINC_AGGREGATION_TERMS_SIZE,default=1000"`

	S3     s3
	MinIO  minIO
	Plugin plugin
}

type s3 struct {
	Bucket string `env:"ZINC_S3_BUCKET"`
}

type minIO struct {
	Endpoint        string `env:"ZINC_MINIO_ENDPOINT"`
	Bucket          string `env:"ZINC_MINIO_BUCKET"`
	AccessKeyID     string `env:"ZINC_MINIO_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"ZINC_MINIO_SECRET_ACCESS_KEY"`
}

type plugin struct {
	ES  elasticsearch
	GSE gse
}

type elasticsearch struct {
	Version string `env:"ZINC_PLUGIN_ES_VERSION"`
}

type gse struct {
	Enable    bool   `env:"ZINC_PLUGIN_GSE_ENABLE,default=false"`
	DictEmbed string `env:"ZINC_PLUGIN_GSE_DICT_EMBED,default=small"`
	DictPath  string `env:"ZINC_PLUGIN_GSE_DICT_PATH,default=./plugins/gse/dict"`
}

var Config = new(config)

func init() {
	err := godotenv.Load()
	fmt.Println(err)
	rv := reflect.ValueOf(Config).Elem()
	loadConfig(rv)
}

func loadConfig(rv reflect.Value) {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)
		if ft.Type.Kind() == reflect.Ptr {
			loadConfig(fv.Elem())
			continue
		} else if ft.Type.Kind() == reflect.Struct {
			loadConfig(fv)
			continue
		}
		if ft.Tag.Get("env") != "" {
			tag := ft.Tag.Get("env")
			setField(fv, tag)
		}
	}
}

func setField(field reflect.Value, tag string) {
	if tag == "" {
		return
	}
	tagColumn := strings.Split(tag, ",")
	v := os.Getenv(tagColumn[0])
	if v == "" {
		if len(tagColumn) > 1 {
			for _, tv := range tagColumn[1:] {
				if strings.HasPrefix(tv, "default=") {
					v = tv[8:]
				}
			}
		}
	}
	if v == "" {
		return
	}
	switch field.Kind() {
	case reflect.Int:
		vi, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal().Err(err).Msgf("env %s is not int", tag)
		}
		field.SetInt(int64(vi))
	case reflect.Bool:
		vi, err := strconv.ParseBool(v)
		if err != nil {
			log.Fatal().Err(err).Msgf("env %s is not bool", tag)
		}
		field.SetBool(vi)
	case reflect.String:
		field.SetString(v)
	default:
		// noop
	}
}
