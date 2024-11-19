package utils

import (
	"os"
)

type FileTask struct {
	SourcePath string
	FileName   string
}

func MoveFile(sourcePath, destPath string) error {
	return os.Rename(sourcePath, destPath)
}
