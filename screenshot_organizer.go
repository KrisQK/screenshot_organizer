package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// 定义源文件夹和目标文件夹路径
	homeDir, _ := os.UserHomeDir()
	sourceDir := filepath.Join(homeDir, "Downloads")
	baseDestDir := filepath.Join(homeDir, "Pictures", "截图整理")

	// 创建目标基础文件夹（如果不存在）
	if err := os.MkdirAll(baseDestDir, 0755); err != nil {
		fmt.Printf("创建目标文件夹失败: %v\n", err)
		return
	}

	// 创建新的截图文件夹（使用当前时间）
	now := time.Now()
	newFolderName := fmt.Sprintf("截图_%s", now.Format("2006-01-02_150405"))
	destDir := filepath.Join(baseDestDir, newFolderName)

	if err := os.Mkdir(destDir, 0755); err != nil {
		fmt.Printf("创建截图文件夹失败: %v\n", err)
		return
	}

	// 读取下载文件夹中的文件
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		fmt.Printf("读取下载文件夹失败: %v\n", err)
		return
	}

	// 移动文件
	movedCount := 0
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "Jietu") {
			sourcePath := filepath.Join(sourceDir, file.Name())
			destPath := filepath.Join(destDir, file.Name())

			// 移动文件
			if err := moveFile(sourcePath, destPath); err != nil {
				fmt.Printf("移动文件 %s 失败: %v\n", file.Name(), err)
				continue
			}
			movedCount++
		}
	}

	fmt.Printf("整理完成！共移动 %d 个文件到 %s\n", movedCount, destDir)
}

func moveFile(sourcePath, destPath string) error {
	// 先尝试直接移动
	if err := os.Rename(sourcePath, destPath); err == nil {
		return nil
	}

	// 如果直接移动失败，则复制后删除
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, source); err != nil {
		return err
	}

	return os.Remove(sourcePath)
}
