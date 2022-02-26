package core

import (
	"context"
	"math"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/blugelabs/bluge"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"gopkg.in/segmentio/analytics-go.v3"

	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

// Telemetry instance
var Telemetry = newTelemetry()

type telemetry struct {
	instanceID   string
	events       chan analytics.Track
	baseInfo     map[string]interface{}
	baseInfoOnce sync.Once
}

func newTelemetry() *telemetry {
	t := new(telemetry)
	t.events = make(chan analytics.Track, 100)
	t.initBaseInfo()

	go t.runEvents()

	return t
}

func (t *telemetry) createInstanceID() string {
	instanceID := uuid.New().String()
	data := map[string]interface{}{
		"_id":   "instance_id",
		"value": instanceID,
	}

	metaIndex := ZINC_SYSTEM_INDEX_LIST["_metadata"]
	doc, _ := metaIndex.BuildBlugeDocumentFromJSON("instance_id", &data)
	doc.AddField(bluge.NewTextField("value", instanceID).StoreValue())
	metaIndex.Writer.Update(doc.ID(), doc)

	return instanceID
}

func (t *telemetry) getInstanceID() string {
	if t.instanceID != "" {
		return t.instanceID
	}

	query := bluge.NewTermQuery("instance_id").SetField("_id")
	searchRequest := bluge.NewTopNSearch(1, query)
	reader, _ := ZINC_SYSTEM_INDEX_LIST["_metadata"].Writer.Reader()
	dmi, err := reader.Search(context.Background(), searchRequest)
	if err != nil {
		log.Printf("core.Telemetry.GetInstanceID: error executing search: %v", err)
	}

	next, err := dmi.Next()
	if err == nil && next != nil {
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			if field == "value" {
				t.instanceID = string(value)
			}
			return true
		})
		if err != nil {
			log.Printf("core.Telemetry.GetInstanceID: error accessing stored fields: %v", err)
		}
	}

	if t.instanceID == "" {
		t.instanceID = t.createInstanceID()
	}

	return t.instanceID
}

func (t *telemetry) initBaseInfo() {
	t.baseInfoOnce.Do(func() {
		m, _ := mem.VirtualMemory()
		cpu_count, _ := cpu.Counts(true)
		zone, _ := time.Now().Local().Zone()

		t.baseInfo = map[string]interface{}{
			"os":           runtime.GOOS,
			"arch":         runtime.GOARCH,
			"zinc_version": v1.Version,
			"time_zone":    zone,
			"cpu_count":    cpu_count,
			"total_memory": m.Total / 1024 / 1024,
		}
	})
}

func (t *telemetry) Instance() {
	if zutils.GetEnv("ZINC_TELEMETRY", "enabled") == "disabled" {
		return
	}

	traits := analytics.NewTraits().
		Set("index_count", len(ZINC_INDEX_LIST)).
		Set("total_index_size_mb", t.TotalIndexSize())

	for k, v := range t.baseInfo {
		traits.Set(k, v)
	}

	v1.SEGMENT_CLIENT.Enqueue(analytics.Identify{
		UserId: t.getInstanceID(),
		Traits: traits,
	})
}

func (t *telemetry) Event(event string, data map[string]interface{}) {
	if zutils.GetEnv("ZINC_TELEMETRY", "enabled") == "disabled" {
		return
	}

	props := analytics.NewProperties()
	for k, v := range t.baseInfo {
		props.Set(k, v)
	}
	for k, v := range data {
		props.Set(k, v)
	}

	t.events <- analytics.Track{
		UserId:     t.getInstanceID(),
		Event:      event,
		Properties: props,
	}
}

func (t *telemetry) runEvents() {
	for event := range t.events {
		v1.SEGMENT_CLIENT.Enqueue(event)
	}
}

func (t *telemetry) TotalIndexSize() float64 {
	TotalIndexSize := 0.0
	for k := range ZINC_INDEX_LIST {
		TotalIndexSize += t.GetIndexSize(k)
	}
	return math.Round(TotalIndexSize)
}

func (t *telemetry) GetIndexSize(indexName string) float64 {
	size := 0.0
	indexType := ZINC_INDEX_LIST[indexName].IndexType

	switch indexType {
	case "s3":
		return size // TODO: implement later
	case "minio":
		return size // TODO: implement later
	default:
		path := zutils.GetEnv("ZINC_DATA_PATH", "./data")
		indexLocation := filepath.Join(path, indexName)
		size, _ = zutils.DirSize(indexLocation)
		return math.Round(size)
	}
}

func (t *telemetry) HeartBeat() {
	m, _ := mem.VirtualMemory()
	data := make(map[string]interface{})
	data["index_count"] = len(ZINC_INDEX_LIST)
	data["total_index_size_mb"] = t.TotalIndexSize()
	data["memory_used_percent"] = m.UsedPercent
	t.Event("heartbeat", data)
}

func (t *telemetry) Cron() {
	c := cron.New()

	c.AddFunc("@every 30m", t.HeartBeat)
	c.Start()
}
