package durability

import (
	"fmt"
)

func WriteWAL(message []byte) {
	walIndexCounter.Inc()
	val := walIndexCounter.GetValue()

	fmt.Println("val: ", val)

	// write WAL entry to batch. This batch will be flushed to the WAL at specified intervals. Check start.go
	// for the flush interval and d=implemnentation details.
	WALBatch.Write(val, message)
}

// return the requested number of WAL entries
func ReadWAL(entriesCount uint64) [][]byte {
	var entries [][]byte
	fmt.Println("durability package readWAL")

	start, err := log1.FirstIndex()
	if err != nil {
		fmt.Println("err: ", err)
	}
	for i := start; i < entriesCount; i++ {
		fmt.Println("durability package readWAL loop")
		entry, err := log1.Read(i)
		if err != nil {
			fmt.Println("durability package readWAL loop error, err: ", err)
		} else {
			entries = append(entries, entry)
		}
	}
	return entries
}

// TODO: TruncateWAL truncates the the consumed WAL entries
func TruncateWAL() {
	fmt.Println("durability package truncateWAL")
}

//TODO: MarkWALAsConsumed marks the WAL entries as consumed
func MarkWALAsConsumed(entriesCount uint64) {
	fmt.Println("durability package markWALAsConsumed")
}
