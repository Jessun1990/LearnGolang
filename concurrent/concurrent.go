package concurrent

import (
	"fmt"
	"math/rand"
	"time"
)

func goroutineExample() {
	doWork := func(strings <-chan string) <-chan interface{} {
		fmt.Println("doWork start.Go func start")
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings {
				fmt.Println(s)
			}
		}()
		fmt.Println("doWork start.Go func end")
		return completed
	}

	doWork(nil)
	fmt.Println("Done.")
}

func goroutineExample2() {
	doWork := func(done <-chan interface{}, strs <-chan interface{}) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strs:
					// do something
					fmt.Println(s)
					fmt.Println("do something")
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
		fmt.Println("Canceling doWork goroutine ...")
		close(done)
	}()
	<-terminated
	fmt.Println("Done.")
}

func goroutineExample3() {
	newRandStream := func() <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.") // 永远不会被执行，永远卡在 for 语句里
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

func goroutineExample4() {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited")
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
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)
	time.Sleep(time.Second)
}

func goroutineExample5() {
	or := func(chans ...<-chan interface{}) <-chan interface{} {
		switch len(chans) {
		case 0:
			return nil
		case 1:
			return chans[0]
		}

		orDone := make(chan interface{})

		go func() {
			defer close(orDone)
			switch len(chans) {
			case 2:
				select {
				case <-chans[0]:
				case <-chans[1]:
				}
			default:
				select {
				case <-chans[0]:
				case <-chans[1]:
				case <-chans[2]:
				case <-or(append(chans[3:], orDone)...):
				}
			}
		}()
	}
}
