package core

import (
	"fmt"
	"time"

	"github.com/blugelabs/bluge"
	blugeindex "github.com/blugelabs/bluge/index"
	"github.com/rs/zerolog/log"
)

const walMaxBatched = 100

// StartWALConsumer Starts a goroutine that consumes the log in a fixed interval
func (index *Index) StartWALConsumer(recover bool) error {
	if err := index.wal.Init(index.Name); err != nil {
		return err
	}

	go func() {
		// Recover & empty WAL on startup
		if recover {
			if err := consumeWAL(index, index.wal.Recover()); err != nil {
				log.Panic().Msgf("WAL recovery of index %v failed", index.Name)
			}
		}

		for {
			// TODO: make configurable
			// TODO: Account for file size as well
			time.Sleep(1 * time.Second)
			err := consumeWAL(index, index.wal.Consume())
			if err != nil {
				log.Warn().Msgf("WAL consumtion for %v failed", index.Name)
			}
		}
	}()

	return nil
}

// consumeWAL Applies all operations from a LogReader in batches to the current index
func consumeWAL(index *Index, lr LogReader) error {
	writers, err := index.GetWriters()
	if err != nil {
		return err
	}

	batched := 0
	batches := make([]*blugeindex.Batch, len(writers))
	for i := range batches {
		batches[i] = bluge.NewBatch()
	}

	flushBatches := func() error {
		if batched <= 0 {
			return nil
		}
		fmt.Println("flush", batched)
		for i := range batches {
			if err := writers[i].Batch(batches[i]); err != nil {
				return err
			}
			batches[i] = bluge.NewBatch()
		}
		batched = 0
		return nil
	}

	for lr.HasNext() {
		doc := lr.Next()

		if doc.Op == LogOpDelete {
			for _, batch := range batches {
				batch.Delete(bluge.NewDocument(doc.DocID).ID())
			}
		} else {
			bdoc, err := index.BuildBlugeDocumentFromJSON(doc.DocID, doc.Doc)
			if err != nil {
				return err
			}

			writer, err := index.FindID(doc.DocID)

			if err == errIdNotFound {
				batches[len(batches)-1].Insert(bdoc)
			} else if doc.Op == LogOpUpdate {
				for i := range writers {
					if writers[i] == writer {
						batches[i].Update(bdoc.ID(), bdoc)
						break
					}
				}
			}
		}
		batched += 1

		if batched >= walMaxBatched {
			if err := flushBatches(); err != nil {
				return err
			}
		}
	}

	if err := flushBatches(); err != nil {
		return err
	}

	return lr.Close()
	//return nil
}
