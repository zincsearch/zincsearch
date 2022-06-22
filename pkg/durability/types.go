package durability

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/zinclabs/zinc/pkg/config"
)

type WALCounter struct {
	mu         sync.Mutex
	count      int
	muMarker   sync.Mutex
	markerFile *os.File
	marker     uint64
}

func (c *WALCounter) OpenMarker(filename string) {
	markerLocation := filepath.Join(config.Global.DataPath, "_wal", filename)
	c.markerFile, _ = os.OpenFile(markerLocation, os.O_RDWR|os.O_CREATE, 0666)
}

// Set the counter to the specified value
func (c *WALCounter) SetMarkerValue(value uint64) {
	c.muMarker.Lock()
	defer c.muMarker.Unlock()
	c.marker = value
	c.markerFile.Write([]byte(fmt.Sprintf("%d", value)))

}

// Increment the counter by 1
func (c *WALCounter) IncMarker() {
	c.muMarker.Lock()
	defer c.muMarker.Unlock()
	c.marker++

}

// TODO: implement this
func (c *WALCounter) GetMarkerValue() uint64 {
	c.muMarker.Lock()
	defer c.muMarker.Unlock()
	val, _ := c.markerFile.Read(make([]byte, 10))
	val, _ = strconv.Atoi(string(rune(val)))
	c.marker = uint64(val)

	return c.marker
}

// Set the counter to the specified value
func (c *WALCounter) SetValue(value int) {
	c.mu.Lock()
	c.count = value
	c.mu.Unlock()
}

// Increment the counter by 1
func (c *WALCounter) Inc() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

func (c *WALCounter) GetValue() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return uint64(c.count)
}
