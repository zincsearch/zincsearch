package test

import (
	"os"
	"testing"

	"github.com/prabhatsharma/zinc/pkg/handlers"
)

func BenchmarkBulk(b *testing.B) {
	f, err := os.Open("../tmp/olympics.ndjson")
	if err != nil {
		b.Error(err)
	}

	target := "olympics"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = handlers.BulkHandlerWorker(target, f)
		if err != nil {
			b.Error(err)
		}
	}
}
