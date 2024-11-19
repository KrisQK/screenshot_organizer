package app

import (
	"fmt"
	"os"
	"path/filepath"
	"screenshot-organizer/internal/config"
	"screenshot-organizer/internal/utils"
	"screenshot-organizer/pkg/ui"
	"strings"
	"sync"
	"time"
)

type App struct {
	cfg *config.Config
}

func New(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Run() error {
	// 显示欢迎界面
	ui.ShowWelcome()

	// 创建目标文件夹
	destDir, err := a.createDestinationDir()
	if err != nil {
		return fmt.Errorf("创建目标文件夹失败: %v", err)
	}

	// 获取需要移动的文件
	tasks, err := a.findScreenshots()
	if err != nil {
		return fmt.Errorf("查找截图文件失败: %v", err)
	}

	if len(tasks) == 0 {
		fmt.Println("⚠️ 没有找到需要移动的截图文件")
		return nil
	}

	// 移动文件
	if err := a.moveFiles(tasks, destDir); err != nil {
		return fmt.Errorf("移动文件失败: %v", err)
	}

	return nil
}

func (a *App) createDestinationDir() (string, error) {
	// 创建基础目录
	if err := os.MkdirAll(a.cfg.BaseDestDir, 0755); err != nil {
		return "", err
	}

	// 创建带时间戳的目标目录
	now := time.Now()
	newFolderName := fmt.Sprintf("截图_%s", now.Format("2006-01-02_150405"))
	destDir := filepath.Join(a.cfg.BaseDestDir, newFolderName)

	if err := os.Mkdir(destDir, 0755); err != nil {
		return "", err
	}

	fmt.Printf("✅ 创建目标文件夹成功: %s\n", destDir)
	return destDir, nil
}

func (a *App) findScreenshots() ([]utils.FileTask, error) {
	files, err := os.ReadDir(a.cfg.SourceDir)
	if err != nil {
		return nil, err
	}

	var tasks []utils.FileTask
	for _, file := range files {
		if strings.HasPrefix(file.Name(), a.cfg.FileNamePrefix) {
			tasks = append(tasks, utils.FileTask{
				SourcePath: filepath.Join(a.cfg.SourceDir, file.Name()),
				FileName:   file.Name(),
			})
		}
	}

	fmt.Printf("🔍 找到 %d 个截图文件需要移动\n", len(tasks))
	return tasks, nil
}

func (a *App) moveFiles(tasks []utils.FileTask, destDir string) error {
	bar := ui.NewProgressBar(len(tasks))

	var wg sync.WaitGroup
	taskChan := make(chan utils.FileTask, len(tasks))
	errChan := make(chan error, len(tasks))

	// 启动工作协程
	for i := 0; i < a.cfg.WorkerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				destPath := filepath.Join(destDir, task.FileName)
				if err := utils.MoveFile(task.SourcePath, destPath); err != nil {
					errChan <- fmt.Errorf("移动文件 %s 失败: %v", task.FileName, err)
				}
				bar.Add(1)
			}
		}()
	}

	// 发送任务
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	// 等待完成
	wg.Wait()
	close(errChan)

	// 处理错误
	for err := range errChan {
		fmt.Printf("\n❌ %v\n", err)
	}

	fmt.Printf("\n✨ 整理完成！共移动 %d 个文件到:\n%s\n", len(tasks), destDir)
	return nil
}
