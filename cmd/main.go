package main

import (
	"log"
	"os"
	"path/filepath"
	"screenshot-organizer/internal/app"
	"screenshot-organizer/internal/config"
)

func main() {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("无法获取用户主目录:", err)
	}

	// 设置源目录和目标目录
	cfg := &config.Config{
		SourceDir:      filepath.Join(homeDir, "Downloads"),        // 源目录设置为下载文件夹
		BaseDestDir:    filepath.Join(homeDir, "Pictures", "截图整理"), // 目标目录设置为图片文件夹下的截图整理
		FileNamePrefix: "Jietu",                                    // 截图文件名前缀
		WorkerCount:    4,                                          // 工作协程数量
	}

	app := app.New(cfg)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
