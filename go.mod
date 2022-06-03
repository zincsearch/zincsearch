module github.com/zinclabs/zinc

go 1.16

require (
	github.com/aws/aws-sdk-go-v2/config v1.11.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.22.0
	github.com/blugelabs/bluge v0.1.9
	github.com/blugelabs/bluge_segment_api v0.2.0
	github.com/blugelabs/query_string v0.3.0
	github.com/bwmarrin/snowflake v0.3.0
	github.com/getsentry/sentry-go v0.13.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.7
	github.com/go-ego/gse v0.70.0
	github.com/goccy/go-json v0.9.6
	github.com/joho/godotenv v1.4.0
	github.com/minio/minio-go/v7 v7.0.21
	github.com/robfig/cron/v3 v3.0.0
	github.com/rs/zerolog v1.26.1
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/stretchr/testify v1.7.1
	github.com/zsais/go-gin-prometheus v0.0.0-20200217150448-2199a42d96c1
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce
	golang.org/x/text v0.3.7
	gopkg.in/segmentio/analytics-go.v3 v3.1.0
)

require (
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/dgraph-io/badger/v3 v3.2103.2
	github.com/google/uuid v1.3.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/prometheus/client_golang v1.12.0 // indirect
	github.com/segmentio/backo-go v1.0.0 // indirect
	github.com/swaggo/files v0.0.0-20210815190702-a29dd2bc99b2
	github.com/swaggo/gin-swagger v1.4.3
	github.com/swaggo/swag v1.8.2
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.etcd.io/etcd/client/v3 v3.5.4
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
)

replace github.com/blugelabs/bluge => github.com/zinclabs/bluge v1.1.0

replace github.com/blugelabs/ice => github.com/zinclabs/ice v1.1.1

replace github.com/blugelabs/bluge_segment_api => github.com/zinclabs/bluge_segment_api v1.0.0
