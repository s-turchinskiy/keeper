package filecheckerutils

import (
	"fmt"
	"os"
)

type FileChecker struct {
	statFunc func(string) (os.FileInfo, error)
}

func NewFileChecker() *FileChecker {
	return &FileChecker{
		statFunc: os.Stat,
	}
}

func (f *FileChecker) CheckFileSize(filePath string, maxFileSize int64) error {
	info, err := f.statFunc(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	if info.Size() > maxFileSize {
		return fmt.Errorf("file too large: %d bytes (max: %d bytes)",
			info.Size(), maxFileSize)
	}

	return nil
}
