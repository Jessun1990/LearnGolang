// Package chapter4 Go 语言的并发模式
// 本书重点，主要是各种并发模式的不同组合
package chapter4

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// codeExample
// "特定约束"举例
// loopData 函数和handleData channel 上的循环都可以使用整数的数据切片
// 惯例是，使用 loopData 函数来访问
func codeExample() {
	data := make([]int, 4)

	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}

// codeExample2
// 通过 channel 的读和写来控制并发进程
func codeExample2() {
	chanOwner := func() <-chan int {
		// 在函数闭包内包含 channel 的写入处理
		// 防止其他 goroutine 写入
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int) { // 内部约束为只读 channel
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}

	results := chanOwner()
	consumer(results)
}

// codeExample3
// 并不是并发安全的数据结构的约束的例子
// 因为传递的切片的不同，在词法范围的原因，已经不可能执行错误的操作。
// 所以不需要通过通信完成内存访问同步或数据共享
func codeExample3() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])

	wg.Wait()
}

// forSelectExample
// for-select 循环
func forSelectExample() {
	for { // 无限循环，使用 range 循环
		select {
		// 使用 channel 进行作业
		}
	}
}

// foreSelectExample2
// 向 channel 发送迭代变量
func forSelectExample2() {
	var done chan interface{}
	var stringStream chan string
	// ...
	for _, s := range []string{"a", "b", "c"} {
		select {
		case <-done:
			return
		case stringStream <- s:
		}
	}
}

// forSelectExample3
// 循环等待停止，变体1
func forSelectExample3() {
	var done chan interface{}
	for {
		select {
		case <-done:
			return
		default:
		}
		// 进行非抢占式任务
	}
}

// forSelectExample4
// 循环等待停止，变体2
func forSelectExample4() {
	var done chan interface{}

	for {
		select {
		case <-done:
			return
		default:
			// 进行非抢占式任务
		}
	}
}

/*
	防止 goroutine 泄漏
*/

// goroutine 有以下几种方式被终止：
// 当它完成了它的工作
// 因为不可恢复的错误，它不能继续工作。
// 当它被告知需要终止工作。

// goroutineExample
// 简单演示 goroutine 泄漏
func goroutineExample() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings {
				// 其他业务代码
				fmt.Println(s)
			}
		}()
		return completed
	}
	doWork(nil)
	// 由于空 goroutine 传递给了 work ，因此字符串永远不会写入任何值
	// 并且包含 doWork 的 goroutine 将在此过程的整个生命周期中保留在内存中
	// 如果在 doWork 和 main goroutine 中加入了 goroutine，甚至会死锁。
	// 解决办法：
	// 在父 goroutine 和其子 goroutine 之间建立一个信号，让父 goroutine 向
	// 其子 goroutine 发出信号通知。

	// 其他操作
	fmt.Println("Done.")
}

// forSelectExample5
// 增加 done 的只读 channel
// 在 channel 上接受 goroutine 的情况
func forSelectExample5() {
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)

			for {
				select {
				case s := <-strings:
					// 业务代码
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

	go func() { // 创建的另一个 goroutine，目的就是为了 close(done)
		time.Sleep(time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done) // 取消的 doWork 中的 goroutine
	}()

	<-terminated // 阻塞直到 close(done)
	fmt.Println("Done.")
}

// forSelectExample6
// 一个 goroutine 阻塞了向 channel 的写入请求
func forSelectExample6() {
	newRandStream := func() <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.") // 这句永远用不会运行
			defer close(randStream)
			for {
				randStream <- rand.Int()
				// 因为无法退出，所以永远阻塞在这里
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

// goroutineExample4
// 在上面的示例中，增加了 done <- chan interface{} 来退出 goroutine
// 如果 goroutine 负责创建 goroutine，那么也要负责确保它可以停止 goroutine
func forSelectExample7() {
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
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)
	time.Sleep(time.Second)
}

// orChannelExample
// 通过递归和 goroutine 创建一个复合 done channel
// 可以将任意数量的 channel 组合到单个 channel 中，
// 只要任何组件 channel 关闭或写入，该 channel 就会关闭。
// 本例将经过一段时间后关闭channel，并将这些 channel 合并到一个关闭的单个 channel 中。
func orChannelExample() {
	var or func(channels ...<-chan interface{}) <-chan interface{}

	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) { // 递归函数的终止条件
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() { // 函数主体
			defer close(orDone)
			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[:3], orDone)...):
				}
			}
		}()
		return orDone
	}

	sig := func(after time.Duration) <-chan interface{} {
		// 创建一个channle，等待指定时间后关闭。
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now() // 大致追踪 or 函数的 channel 何时开始阻塞

	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
	)

	fmt.Printf("done after %v\n", time.Since(start))
}

// 输出结果：
//	done afater 1.096730468s
// 尽管在调用中放置了多个 channel 或需要不同时间才能关闭
// 但是 1s 后，关闭的那个channel 会导致整个 channel 关闭

/*
	错误处理
*/
// errHandleExample
// 简单的错误处理示例
func errHandleExample() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan *http.Response {
		responses := make(chan *http.Response)
		go func() {
			defer close(responses)
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}
				select {
				case <-done:
					return
				case responses <- resp:
				}
			}
		}()
		return responses
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.baidu.com", "https://badhost"}
	for response := range checkStatus(done, urls...) {
		fmt.Printf("Response: %v\n", response.Status)
	}
}

//errHandleExample2
// 上面示例的更佳的解决方案
func errHandleExample2() {
	type Result struct {
		Error    error
		Response *http.Response
	}

	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				result = Result{Error: err, Response: resp}
				select {
				case <-done:
					return
				case results <- result: // 请求的结果返回
				}
			}
		}()
		return results
	}

	done := make(chan interface{})
	defer close(done)
	urls := []string{"https://www.baidu.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v", result.Error)
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

// errHandleExample3
// 上面示例的修改版
func errHandleExample3() {
	done := make(chan interface{})
	defer close(done)

	type Result struct {
		Error    error
		Response *http.Response
	}

	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				result = Result{Error: err, Response: resp}
				select {
				case <-done:
					return
				case results <- result: // 请求的结果返回
				}
			}
		}()
		return results
	}

	errCount := 0
	urls := []string{"a", "https://www.baidu.com", "b", "c", "d"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			errCount++
			if errCount >= 3 { // 错误超过3个时，跳出 range 循环
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

// 在构建 goroutine 的返回值时，应将错误视为一等公民。
// 如果你的 goroutine 可能产生错误，那么这些错误应该与你的结果类型紧密结合
// 并且通过相同的通信线传递，就像常规的同步函数。

/*
	pipeline
*/
// 需要流式处理或批处理数据时， pipeline 是一个非常强大的工具
// pipeline 只不过是一系列将数据输入，执行操作并将结果传回的系统。这些操作称之为 stage
// 通过 pipeline，可以分离每个 stage 的关注点。

// piplineExample
// 简单的 pipline stage 示例
func pipelineExample() {
	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}
		return multipliedValues
	}

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}

	// 尝试将上面两个函数合并
	ints := []int{1, 2, 3, 4}
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}
	// 上面 stage 执行的是批处理操作，stage 一次只接收和处理一个元素，则是流处理。
	// 为了保持原始数据不变，每个 stage 必须创建一个等长的信片段来存储其计算结果。
}

// piplineExample2
// 每个 stage 都接收并发出一个离散值，内存占用回落到只有 pipeline 输入的大小。
// 但是 pipeline 写入到 for-range 内部，限制了 pipline 的重复使用。
func pipelineExample2() {
	multiply := func(value, multiplier int) int {
		return value * multiplier
	}

	add := func(value, additive int) int {
		return value + additive
	}

	ints := []int{1, 2, 3, 4}
	for _, v := range ints {
		fmt.Println(multiply(add(multiply(v, 2), 1), 2))
	}
}

/*
	构建 pipeline 的最佳实践
*/

// pipelineExample3
// 创建一个 done channel，并在 defer 中关闭。保证 goroutine 不泄漏。
func pipelineExample3() {
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		// 参数：接收一个可变的整数切片
		// 使用 goroutine 来返回切片整数的 channel
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multiply := func(
		done <-chan interface{},
		intStream <-chan int,
		multiplier int) <-chan int {
		multipliedStram := make(chan int)
		go func() {
			defer close(multipliedStram)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStram <- i * multiplier:
				}
			}
		}()
		return multipliedStram
	}

	add := func(
		done <-chan interface{},
		intStream <-chan int,
		additive int,
	) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- i + additive:
				}
			}
		}()
		return addedStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)

	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline { // 对传入的 channel 进行迭代，当 channle 被 generator 关闭，
		// range 语句也就结束了。
		fmt.Println(v)
	}

	// 在 pipeline 开始时，已经确定将离散值转换为 channel，这一过程中，两点必须是可抢占的:
	// 创建几乎不是瞬时的离散值
	// 在 channel 上发送离散值
}

/*
	一些便利的生成器
*/

// generatorExample
func generatorExample() {
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	take := func( // 只会从其传入的 valueStream 中取出第一个 num 项目
		done <-chan interface{},
		valueStream <-chan interface{},
		num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- valueStream:
				}
			}
		}()
		return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	for num := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%+v ", num)
	}
}

// generatorExample2
// 重复调用函数的生成器
func generatorExample2() {
	take := func( // 只会从其传入的 valueStream 中取出第一个 num 项目
		done <-chan interface{},
		valueStream <-chan interface{},
		num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- valueStream:
				}
			}
		}()
		return takeStream
	}

	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}

	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} { return rand.Int() }

	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}
	// 根据需要生成随机整数的无限 channel
}

// forSelectExample8
func forSelectExample8() {
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	take := func( // 只会从其传入的 valueStream 中取出第一个 num 项目
		done <-chan interface{},
		valueStream <-chan interface{},
		num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- valueStream:
				}
			}
		}()
		return takeStream
	}
	toString := func(done <-chan interface{},
		valueStream <-chan interface{}) <-chan string {
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case stringStream <- v.(string):
				}
			}
		}()
		return stringStream
	}

	done := make(chan interface{})
	defer close(done)

	var message string
	for token := range toString(done, take(done, repeat(done, "I", "am."), 5)) {
		message += token
	}
	fmt.Printf("message: %s...", message)
}

// TODO

/*
	扇入，扇出/ Fan-in, Fan-out
*/
// 扇出，用于描述启动多个 goroutine 以处理来自 pipline 的输入的过程，
// 扇入，是将多个结果组合到一个 channel 的过程。
// TODO

func fanInFanOutExmaple() {
}

/*
	or-done-channel
*/
func orDoneExample() {
	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		// 使用 orDone 来封装多层嵌套的循环
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if ok == false {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	var done chan interface{}
	var myChan chan interface{}
	// ...
	// ...

	for val := range orDone(done, myChan) {
		fmt.Println(val)
	}
}
