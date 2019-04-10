package channel

import (
	"fmt"
	"time"
)

// TryGoroutine1 ...
func TryGoroutine1() {
	doWork := func(strings <-chan string) <-chan interface{} {

		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings {
				fmt.Println(s)
			}
		}()
		return completed
	}

	res := doWork(nil)
	fmt.Printf("Done. res: %+v", res)
}

// TryGoroutine2 ...
func TryGoroutine2() {
	doWork := func(done <-chan interface{},
		strings <-chan string) <-chan interface{} {

		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done.")
}
