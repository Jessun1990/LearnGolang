package channel

import (
	"fmt"
	"math/rand"
	"time"
)

// TryGoroutine1 ...
// 由于 doWork 传递了 nil，所以 doWork 中的 goroutine 会泄漏.
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
// 加入了 done ， 可以消除 goroutine 的泄漏
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

// TryGoroutine3 ...
func TryGoroutine3() {
	newRandStream := func() <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				randStream <- rand.Int()
			}
		}()
		return randStream
	}

	randStream := newRandStream()
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
}

// TryGoroutine4 ...
func TryGoroutine4() {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randStream
	}

	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, randStream)
	}

	close(done)
}
