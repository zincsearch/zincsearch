package zutils

import (
	"os"
	"path/filepath"
)

// DirSize return the size of the directory (KB)
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
	sizeKB := float64(size) / 1024.0

	return sizeKB, err
}

func IsExist(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsExist(err) {
			return false, nil
		}
		return false, err
	}
	f.Close()

	return true, nil
}
