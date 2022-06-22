package durability

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tidwall/wal"
	"github.com/zinclabs/zinc/pkg/config"
)

var WALBatch *wal.Batch
var timer1 *time.Timer
var log1 *wal.Log

var walIndexCounter WALCounter

func init() {
	fmt.Println("durability package init")
	WALBatch = new(wal.Batch)

	os.MkdirAll(filepath.Join("./data", "_wal"), 0755)

	walLocation := filepath.Join(config.Global.DataPath, "_wal", "logs")

	// Create the marker file of WAL. This file will store the last index of WAL entry that has been consumed.
	consumedWalMarkerLocation := filepath.Join(config.Global.DataPath, "_wal", "marker.txt")

	markerFile, err := os.OpenFile(consumedWalMarkerLocation, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0755)
	if err != nil {
		panic(err)
	}

	data := []byte("9\n")
	markerFile.Write(data)
	markerFile.Seek(0, 0)

	log1, _ = wal.Open(walLocation, nil)

	lastIndex, _ := log1.LastIndex()
	walIndexCounter.SetValue(int(lastIndex))

	timerTime := time.Duration(config.Global.WALSyncTime) * time.Second // default is 5 seconds

	// All the WAL entries will be inserted in a batch
	// set up the timer in a separate goroutine to exeute after specified seconds
	go func() {
		fmt.Println("pkg init go func")
		timer1 = time.NewTimer(timerTime)
		for { // loop forever to create an infinite looping timer
			<-timer1.C
			fmt.Println("timer1 at: ", time.Now())
			timer1.Reset(timerTime) // reset timer to run again after specified seconds
			log1.WriteBatch(WALBatch)
			lastIndex, _ := log1.LastIndex()         // get the last index of the WAL
			walIndexCounter.SetValue(int(lastIndex)) // set the counter to the last index. This can be used
			WALBatch.Clear()
			log1.Sync() // sync WAL to file system
		}
	}()
}
