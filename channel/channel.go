package channel

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

/*
 channel declaration
 var dataStream chan interface{}
 dataStream = make(chan interface{})

 receive-only chan
 var dataStream <-chan interface{}
 dataStream := make(<-chan interface{})

 send-only chan
 var dataStream chan<- interface{}
 dateStream := make(chan<- interface{})
*/

// 通过通信共享内存来进行同步
// channel 是阻塞的
func chanExample1() {
	stringStream := make(chan string)
	go func() {
		//defer close(stringStream)
		stringStream <- "Hello Golang Channels"
	}()
	fmt.Println(<-stringStream)
}

// channel 可以被遍历， channel 必须被 close
func chanExample2() {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for i := 1; i <= 5; i++ {
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Printf("%+v\n", integer)
	}
}

func chanExample3() {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			fmt.Printf("%+v has begun\n", i)
		}(i)
	}
	fmt.Println("Unblocking goroutines...")
	close(begin) // chan 被 close 的时候，goroutines 开始执行
	wg.Wait()
}

/*
	1. 实例化 channel
	2. 执行写操作，或者将所有权传递给另一个 goroutine
	3. 关闭 channel
	4. 确认前三件事，并通过一个只读 channel 将他们暴露
*/

func chanExample4() {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %+v.\n", i)
			intStream <- i
		}

	}()
	for integar := range intStream {
		fmt.Fprintf(&stdoutBuff, "Received: %+v.\n", integar)
	}
}

/*
	缓冲的 channel
*/

func chanExample5() {
	chanOwner := func() <-chan int {
		resultStream := make(chan int, 3)
		go func() {
			defer close(resultStream)
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream
	}

	resultStream := chanOwner()
	for result := range resultStream {
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Println("Done receiving!")
}

/*
	func selectExample1() {
		var c1, c2 <-chan interface{}
		var c3 chan<- interface{}
		select {
		case <-c1:
			//
		case <-c2:
			//
		case c3 <- struct{}{}:
			//
		}
	}
*/

func selectExample1() {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()
	fmt.Println("Blocking on read...")
	select {
	case <-c:
		fmt.Printf("Unblocking %+v later. \n", time.Since(start))
	}
}

func selectExample2() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}
	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}

func selectExample3() {
	//start := time.Now()
	var c <-chan int
	select {
	case <-c:
	case <-time.After(time.Second):
		fmt.Println("Time out.")
		//default:
		//fmt.Printf("In default after %+v\n\n", time.Since(start))
	}
}

func selectExample4() {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
			fmt.Println("This is default")
		}
		workCounter++
		time.Sleep(time.Second)
	}
	fmt.Printf("Achieved %+v cycles of work before signalled to stop. \n", workCounter)
}

// or-channel : 将任意数量的channel组合到单个channel中，只要任何
// 组件 channel 关闭或者写入，该 channel 就会关闭
func selectExample5() {
	var or func(chans ...<-chan interface{}) <-chan interface{}
	or = func(chans ...<-chan interface{}) <-chan interface{} {
		switch len(chans) {
		case 0:
			return nil // 递归函数，终止条件
		case 1:
			return chans[0]
		}

		orDone := make(chan interface{})

		go func() { // 本函数最重要的部分
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

	///////// 例子
	sig := func(after time.Duration) <-chan interface{} { // 创建一个 channel ，指定时间后关闭
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
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %+v", time.Since(start))
}

// 使用 channel 过程中的错误处理
type result struct {
	Error    error
	Response *http.Response
}

func chanExample6() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan result {
		results := make(chan result)
		go func() {
			defer close(results)
			for _, url := range urls {
				var res result
				rsp, err := http.Get(url)
				res = result{
					Error:    err,
					Response: rsp,
				}
				select {
				case <-done:
					return
				case results <- res:
				}
			}
		}()
		return results
	}
	done := make(chan interface{})
	defer close(done)
	urls := []string{"https://www.google.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: <%+v>", result.Error)
			continue
		}
		fmt.Printf("Response: %+v\n", result.Response.Status)
	}
}
