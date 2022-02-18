package zutils

import (
	"os"
	"path/filepath"
)

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
