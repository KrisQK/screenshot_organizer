package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

const asciiArt = `
╭──────────────────────────────────────╮
│     MacOS 截图移动助手 v1.0          │
│     让您的截图管理更轻松 ～          │
╰──────────────────────────────────────╯
`

type FileTask struct {
	sourcePath string
	destPath   string
	fileName   string
}

func main() {
	// 显示欢迎界面
	fmt.Println(asciiArt)
	fmt.Println("程序启动中...")
	time.Sleep(1 * time.Second)

	// 定义源文件夹和目标文件夹路径
	homeDir, _ := os.UserHomeDir()
	sourceDir := filepath.Join(homeDir, "Downloads")
	baseDestDir := filepath.Join(homeDir, "Pictures", "截图整理")

	fmt.Printf("📂 源文件夹: %s\n", sourceDir)
	fmt.Printf("📂 目标文件夹: %s\n", baseDestDir)

	// 创建目标基础文件夹
	if err := os.MkdirAll(baseDestDir, 0755); err != nil {
		fmt.Printf("❌ 创建目标文件夹失败: %v\n", err)
		return
	}

	// 创建新的截图文件夹
	now := time.Now()
	newFolderName := fmt.Sprintf("截图_%s", now.Format("2006-01-02_150405"))
	destDir := filepath.Join(baseDestDir, newFolderName)

	if err := os.Mkdir(destDir, 0755); err != nil {
		fmt.Printf("❌ 创建截图文件夹失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 创建目标文件夹成功: %s\n", destDir)

	// 读取下载文件夹中的文件
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		fmt.Printf("❌ 读取下载文件夹失败: %v\n", err)
		return
	}

	// 筛选需要移动的文件
	var tasks []FileTask
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "Jietu") {
			tasks = append(tasks, FileTask{
				sourcePath: filepath.Join(sourceDir, file.Name()),
				destPath:   filepath.Join(destDir, file.Name()),
				fileName:   file.Name(),
			})
		}
	}

	if len(tasks) == 0 {
		fmt.Println("⚠️ 没有找到需要移动的截图文件")
		return
	}

	fmt.Printf("🔍 找到 %d 个截图文件需要移动\n", len(tasks))
	time.Sleep(1 * time.Second)

	// 创建进度条
	bar := progressbar.NewOptions(len(tasks),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("📦 正在移动文件..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	// 使用 worker pool 处理文件移动
	var wg sync.WaitGroup
	taskChan := make(chan FileTask, len(tasks))
	workerCount := 3 // 设置3个工作协程

	// 启动工作协程
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				if err := moveFile(task.sourcePath, task.destPath); err != nil {
					fmt.Printf("\n❌ 移动文件 %s 失败: %v\n", task.fileName, err)
				}
				bar.Add(1)
			}
		}()
	}

	// 发送任务到通道
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	// 等待所有工作完成
	wg.Wait()

	fmt.Printf("\n✨ 整理完成！共移动 %d 个文件到:\n%s\n", len(tasks), destDir)
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
