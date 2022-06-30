package core

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path"
	"sync"

	"github.com/zinclabs/zinc/pkg/config"
)

type LogOpType = int

const (
	LogOpCreate LogOpType = 1
	LogOpUpdate LogOpType = 2
	LogOpDelete LogOpType = 3
)

// LogOpInfo Entries are stored in log files
type LogOpInfo struct {
	Op    LogOpType
	DocID string
	Doc   map[string]interface{}
	// consecutive ordinals for operations with the same id
	Version int
}

// WriteAheadLog is a simple Write-Ahead Log implementation.
// It allows making changes to an index concurrently and efficiently processing
// them from a worker thread. At any moment all crucial data is
// kept on persistent storage before being written to the index.
// In case of a crash recovery is possible.
type WriteAheadLog struct {
	readFile  *os.File
	writeFile *os.File
	versions  map[string]int
	mu        sync.Mutex

	//
	written int
}

func (tl *WriteAheadLog) Init(index string) (err error) {
	tl.versions = make(map[string]int)
	tl.written = 0

	// TODO: better file placement on non disk stored indices
	readFileName := config.Global.DataPath + "/" + index + "/wal-1"
	writeFileName := config.Global.DataPath + "/" + index + "/wal-2"

	if err := os.MkdirAll(path.Dir(readFileName), 0770); err != nil {
		return err
	}

	tl.readFile, err = os.OpenFile(readFileName, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return
	}
	tl.writeFile, err = os.OpenFile(writeFileName, os.O_CREATE|os.O_RDWR, 0777)
	return
}

// CreateUpdateDocument Create or update a document
func (tl *WriteAheadLog) CreateDocument(id string, doc map[string]interface{}) error {
	return tl.put(LogOpInfo{DocID: id, Doc: doc, Op: LogOpCreate})
}

func (tl *WriteAheadLog) UpdateDocument(id string, doc map[string]interface{}) error {
	return tl.put(LogOpInfo{DocID: id, Doc: doc, Op: LogOpUpdate})
}

// DeleteDocument DeleteDocument a document
func (tl *WriteAheadLog) DeleteDocument(id string) error {
	return tl.put(LogOpInfo{DocID: id, Op: LogOpDelete})
}

// put writes an operation to the log and keeps track of versioning
func (tl *WriteAheadLog) put(doc LogOpInfo) error {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	if storedVersion, ok := tl.versions[doc.DocID]; ok {
		doc.Version = storedVersion + 1
	}
	tl.versions[doc.DocID] = doc.Version

	jsonData, err := json.Marshal(doc) // TODO: error handling? we got this parsed from gin.JSON...
	if err != nil {
		return err
	}

	if _, err = tl.writeFile.Write(jsonData); err != nil {
		return err
	}
	if _, err = tl.writeFile.Write([]byte{'\n'}); err != nil {
		return err
	}

	tl.written += 1

	return tl.writeFile.Sync()
}

// Consume Returns a LogReader for consuming all entries up to the current timepoint
func (tl *WriteAheadLog) Consume() LogReader {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	tl.swap()

	reader := openLogFile(tl.readFile)
	versions := tl.versions
	tl.versions = make(map[string]int)
	tl.written = 0

	return LogReader{
		reader:   reader,
		versions: versions,
	}
}

// Recover Returns a LogReader for recovering the log
func (tl *WriteAheadLog) Recover() LogReader {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	reader1 := openLogFile(tl.readFile)
	reader2 := openLogFile(tl.writeFile)

	// file earliest written to comes first
	if writeTime(tl.readFile) > writeTime(tl.writeFile) {
		reader1, reader2 = reader2, reader1
	}

	versions := make(map[string]int)
	var maxVersion int = 0
	for reader1.HasNext() {
		doc := reader1.Next()
		versions[doc.DocID] = doc.Version
		if doc.Version > maxVersion {
			maxVersion = doc.Version
		}
	}

	reader2.versOffset = maxVersion + 1
	for reader2.HasNext() {
		doc := reader2.Next()
		versions[doc.DocID] = doc.Version
	}

	reader1.reset()
	reader2.reset()

	return LogReader{
		reader: &DualLogFileReader{
			reader1: reader1,
			reader2: reader2,
		},
		versions: versions,
	}
}

func writeTime(f *os.File) int64 {
	info, _ := f.Stat()
	return info.ModTime().UnixNano()
}

// swap Swaps write & read file and seeks to start
func (tl *WriteAheadLog) swap() {
	tl.readFile, tl.writeFile = tl.writeFile, tl.readFile
	tl.readFile.Seek(0, 0)
}

// singleLogFileReader Reads entries from a single log file
// and truncates the file to zero length on close
type singleLogFileReader struct {
	file   *os.File
	reader *bufio.Scanner
	// offset all versions by fixed number (used for chaining files in order)
	versOffset int
}

func openLogFile(file *os.File) singleLogFileReader {
	reader := singleLogFileReader{file: file}
	reader.reset()
	return reader
}

func (sr singleLogFileReader) HasNext() bool {
	return sr.reader.Scan()
}

func (sr singleLogFileReader) Next() LogOpInfo {
	jsonStr := sr.reader.Text()
	var doc LogOpInfo
	json.Unmarshal([]byte(jsonStr), &doc) // TODO: check err
	doc.Version += sr.versOffset
	return doc
}

func (sr singleLogFileReader) Close() error {
	sr.file.Truncate(0)
	sr.file.Sync()
	return nil
}

func (sr *singleLogFileReader) reset() {
	sr.file.Seek(0, 0)
	sr.reader = bufio.NewScanner(sr.file)
}

// logFileReader Read from one or multiple log files
type logFileReader interface {
	io.Closer
	HasNext() bool
	Next() LogOpInfo
}

// DualLogFileReader Combines two logFileReader to read from two files consecutively
type DualLogFileReader struct {
	reader1   logFileReader
	reader2   logFileReader
	readFirst bool
}

func (dr *DualLogFileReader) HasNext() bool {
	if !dr.readFirst {
		if !dr.reader1.HasNext() {
			dr.readFirst = true
			return dr.HasNext()
		}
		return true
	} else {
		return dr.reader2.HasNext()
	}
}

func (dr *DualLogFileReader) Next() LogOpInfo {
	if !dr.readFirst {
		return dr.reader1.Next()
	} else {
		return dr.reader2.Next()
	}
}

func (dr *DualLogFileReader) Close() error {
	dr.reader1.Close()
	dr.reader2.Close()
	return nil
}

// LogReader Allows reading part of the log and deleting it afterwards
type LogReader struct {
	versions map[string]int
	reader   logFileReader
}

func (lr LogReader) HasNext() bool {
	return lr.reader.HasNext()
}

func (lr LogReader) Next() LogOpInfo {
	for {
		doc := lr.reader.Next()
		if storedVersion := lr.versions[doc.DocID]; storedVersion > doc.Version {
			if !lr.reader.HasNext() {
				panic("impossible") // TODO: rewrite better
			}
			continue
		}
		return doc
	}
}

func (lr LogReader) Close() error {
	return lr.reader.Close()
}
