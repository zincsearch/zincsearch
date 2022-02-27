package v1

import (
	"time"

	"gopkg.in/segmentio/analytics-go.v3"
)

var (
	// SEGMENT_CLIENT := analytics.New("hQYncuWEjDJC23MnU6jHXiye5k7qP2PL")
	SEGMENT_CLIENT analytics.Client
)

func init() {
	cf := analytics.Config{
		Interval:  15 * time.Second,
		BatchSize: 100,
		// Endpoint: "http://localhost:8080/api/v1/segment",
	}

	SEGMENT_CLIENT, _ = analytics.NewWithConfig("hQYncuWEjDJC23MnU6jHXiye5k7qP2PL", cf)
}
