package sync

import (
	"fmt"
	"sync"
	"time"
)

// TryWg wg.WaitGroup 的用法示例
// wg.Wait() 后的语句会等到所有的 goroutine 完成才会开始执行
func TryWg() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("1st goroutine sleeping...")
		time.Sleep(1)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("2nd goroutine sleeping...")
		time.Sleep(1)
	}()

	wg.Wait() // 等待所有的 Done() 完成
	fmt.Println("All goroutines complete")
}
