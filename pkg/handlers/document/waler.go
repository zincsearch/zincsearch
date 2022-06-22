package document

import (
	"io"

	"github.com/zinclabs/zinc/pkg/durability"
)

func sendToWAL(messageType string, docID string, indexName string, body *io.ReadCloser) error {
	separator := 0xff
	// Write to WAL
	message, err := io.ReadAll(*body)
	if err != nil {
		return err
	}

	walMessage := append([]byte(messageType), byte(separator))
	walMessage = append(walMessage, docID...) // its legal to append a string to a byte slice
	walMessage = append(walMessage, byte(separator))
	walMessage = append(walMessage, indexName...)
	walMessage = append(walMessage, byte(separator))
	walMessage = append(walMessage, message...)

	durability.WriteWAL([]byte(walMessage))
	return nil
}
func ReadFromWAL(entriesCount uint64) [][]byte {
	return durability.ReadWAL(entriesCount)
}

func TruncateWAL() {
	durability.TruncateWAL()
}

func MarkWALAsConsumed(entriesCount uint64) {
	
	durability.MarkWALAsConsumed(entriesCount)
}
