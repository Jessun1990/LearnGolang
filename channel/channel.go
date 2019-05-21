package channel

import (
	"bytes"
	"fmt"
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
		}
		workCounter++
		time.Sleep(time.Second)
	}
	fmt.Printf("Achieved %+v cycles of work before signalled to stop. \n", workCounter)
}
