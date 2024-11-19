package ui

import (
	"fmt"
)

type ProgressBar struct {
	total   int
	current int
}

func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		total:   total,
		current: 0,
	}
}

func (p *ProgressBar) Add(n int) {
	p.current += n
	p.display()
}

func (p *ProgressBar) display() {
	percentage := float64(p.current) * 100 / float64(p.total)
	fmt.Printf("\r进度: %.2f%% [%d/%d]", percentage, p.current, p.total)
}

func ShowWelcome() {
	art := `
    ╭━━━━━━━━━━━━━━━━━━━━━━━╮
    │     截图整理工具       │
    │    ⊂((・▽・))⊃        │
    │   让文件井然有序       │
    ╰━━━━━━━━━━━━━━━━━━━━━━━╯
         _____   _____
        /    /  /    /
       /    /  /    /
      /____/  /____/
     /    /  /    /
    /    /  /    /
   /____/  /____/

    -----坤坤出品-----
    `
	fmt.Println(art)
}
