package core

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"runtime"

	"github.com/blugelabs/bluge"
	"github.com/google/uuid"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
	"github.com/prabhatsharma/zinc/pkg/zutils"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"gopkg.in/segmentio/analytics-go.v3"
)

func CreateInstanceID() string {

	metaIndex := ZINC_SYSTEM_INDEX_LIST["_metadata"]

	instance_id := uuid.New().String()

	data := map[string]interface{}{
		"_id":   "instance_id",
		"Value": instance_id,
	}

	doc, _ := metaIndex.BuildBlugeDocumentFromJSON("instance_id", &data)

	doc.AddField(bluge.NewTextField("Value", instance_id).StoreValue())

	ZINC_SYSTEM_INDEX_LIST["_metadata"].Writer.Update(doc.ID(), doc)

	m, _ := mem.VirtualMemory()
	cpu_count, _ := cpu.Counts(true)

	if zutils.GetEnv("ZINC_TELEMETRY", "enabled") == "disabled" {
		return instance_id
	}

	v1.SEGMENT_CLIENT.Enqueue(analytics.Identify{
		UserId: instance_id,
		Traits: analytics.NewTraits().
			Set("os", runtime.GOOS).
			Set("arch", runtime.GOARCH).
			Set("zinc_version", v1.Version).
			Set("cpu_count", cpu_count).
			Set("memory", m.Total),
	})

	return instance_id
}

func GetInstanceID() string {

	metaIndex := ZINC_SYSTEM_INDEX_LIST["_metadata"]

	query := bluge.NewTermQuery("instance_id").SetField("_id")
	// query := bluge.NewTermQuery("instance_id")

	searchRequest := bluge.NewTopNSearch(1, query)

	reader, _ := metaIndex.Writer.Reader()

	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Printf("error executing search: %v", err)
	}

	instance_id := ""

	next, err := dmi.Next()
	for err == nil && next != nil {
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			if field == "Value" {
				instance_id = string(value)
				return true
			}
			return true
		})
		if err != nil {
			log.Printf("error accessing stored fields: %v", err)
		}

		next, err = dmi.Next()
	}

	if instance_id == "" {
		instance_id = CreateInstanceID()
	}

	return instance_id
}

func TelemetryInstance() {
	if zutils.GetEnv("ZINC_TELEMETRY", "enabled") == "disabled" {
		return
	}

	m, _ := mem.VirtualMemory()
	cpu_count, _ := cpu.Counts(true)

	v1.SEGMENT_CLIENT.Enqueue(analytics.Identify{
		UserId: GetInstanceID(),
		Traits: analytics.NewTraits().
			Set("index_count", len(ZINC_INDEX_LIST)).
			Set("total_index_size_mb", TotalIndexSize()).
			Set("os", runtime.GOOS).
			Set("arch", runtime.GOARCH).
			Set("zinc_version", v1.Version).
			Set("cpu_count", cpu_count).
			Set("total_memory", m.Total/1024/1024),
	})
}

func TelemetryEvent(event string, data map[string]interface{}) {

	m, _ := mem.VirtualMemory()
	cpu_count, _ := cpu.Counts(true)

	props := analytics.NewProperties().
		Set("index_count", len(ZINC_INDEX_LIST)).
		Set("total_index_size_mb", TotalIndexSize()).
		Set("os", runtime.GOOS).
		Set("arch", runtime.GOARCH).
		Set("zinc_version", v1.Version).
		Set("cpu_count", cpu_count).
		Set("total_memory", m.Total/1024/1024).
		Set("memory_used_percent", m.UsedPercent)

	for k, v := range data {
		props.Set(k, v)
	}

	if zutils.GetEnv("ZINC_TELEMETRY", "enabled") == "disabled" {
		return
	}

	v1.SEGMENT_CLIENT.Enqueue(analytics.Track{
		UserId:     GetInstanceID(),
		Event:      event,
		Properties: props,
	})
}

func DirSize(path string) (float64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	sizeMB := float64(size) / 1024.0 / 1024.0

	return sizeMB, err
}

func TotalIndexSize() float64 {
	TotalIndexSize := 0.0
	for k := range ZINC_INDEX_LIST {
		path := zutils.GetEnv("ZINC_DATA_PATH", "./data")
		indexLocation := filepath.Join(path, k)
		size, _ := DirSize(indexLocation)
		TotalIndexSize += size
	}

	return math.Round(TotalIndexSize)
}

func GetIndexSize(indexName string) float64 {

	size := 0.0

	indexType := ZINC_INDEX_LIST[indexName].IndexType

	if indexType == "s3" {
		return size // TODO: implement later
	} else if indexType == "minio" {
		return size // TODO: implement later
	} else if indexType == "disk" {
		path := zutils.GetEnv("ZINC_DATA_PATH", "./data")
		indexLocation := filepath.Join(path, indexName)
		size, _ = DirSize(indexLocation)
		return math.Round(size)
	}

	return size
}

func TelemetryCron() {
	c := cron.New()

	c.AddFunc("@every 30m", HeartBeat)
	c.Start()
}

func HeartBeat() {
	TelemetryEvent("heartbeat", map[string]interface{}{})
}
