package goroutine

import (
	"fmt"
	"time"
)

// ShowGoroutine : go test ./goroutine -run TestShowGoroutine -v
func ShowGoroutine() {
	fmt.Println("Example Show: <Goroutine>")
	for i := 0; i <= 10; i++ {
		go func(no int) {
			fmt.Printf("i = %+v\n", no)
		}(i)
		time.Sleep(time.Millisecond)
	}
}
