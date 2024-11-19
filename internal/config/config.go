package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	SourceDir      string
	BaseDestDir    string
	WorkerCount    int
	FileNamePrefix string
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &Config{
		SourceDir:      filepath.Join(homeDir, "Downloads"),
		BaseDestDir:    filepath.Join(homeDir, "Pictures", "截图整理"),
		WorkerCount:    3,
		FileNamePrefix: "Jietu",
	}, nil
}
