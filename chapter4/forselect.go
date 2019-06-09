package chapter4

import (
	"fmt"
	"math/rand"
	"time"
)

/*
for { // 无限循环或者使用 range 循环
	select {
		// 使用 channel 进行作业
	}
}
*/

/*
for _, s := range []string{"a", "b", "c"}{
	select {
		case <- done :
			return
		case stringStream <- s:
	}
}
*/

/*
循环等待停止
第一种
for {
	select {
		case <- done:
			return
		default:
	}
	// 进行抢占式任务
}
第二种
for {
	select {
		case <- done:
			return
		default:
			// 进行抢占式任务
	}
}
*/

// 父 goroutine 向子 goroutine 发出信号通知
// 处理 channel 上接收 goroutine 的情况
func goroutineExample() {
	doWork := func(
		done <-chan interface{}, strings <-chan string) <-chan interface{} {

		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
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
		time.Sleep(time.Second)
		fmt.Println("Canceling doWork goroutine...")
	}()

	<-terminated
	fmt.Println("Done.")
}

// goroutine 阻塞了向 channel 进行写入的请求
func goroutineExample2() {
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
		fmt.Printf("%+v: %+v\n", i, <-randStream)
	}
	close(done)
	time.Sleep(time.Second)
}

// 通过递归和 goroutine 创建一个复合的 done channel
func goroutineExample3() {
	var or func(chans ...<-chan interface{}) <-chan interface{}
	or = func(chans ...<-chan interface{}) <-chan interface{} {
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
		return orDone
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(3*time.Minute),
		sig(4*time.Second),
	)
	fmt.Printf("done after %+v", time.Since(start))
}
