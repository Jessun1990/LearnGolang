package chapter3

import (
	"bytes"
	"fmt"
	"os"
	"sync"
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

// chanExample channel 的基本用法展示，channel 是阻塞的
func chanExample() {
	stringStream := make(chan string)
	go func() {
		stringStream <- "Hello World"
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

// chanExample3 : channel 被 close 的时候，所有的 goroutines 开始执行
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
	close(begin)
	wg.Wait()
}

/*
	1. 实例化 channel
	2. 执行写操作，或者将所有权传递给另一个 goroutine
	3. 关闭 channel
	4. 确认前三件事，并通过一个只读 channel 将他们暴露
*/

// chanExample4 带缓冲区的 channel
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

// chanExample5 resultStream 的生命周期封装在 chanOwner 中
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
		fmt.Printf("Received: %+v\n", result)
	}
	fmt.Println("Done receiving!")
}
