// Package chapter3 Go 语言并发组件
// 本书重点，主要关于 gorountine 和 sync 相关内容
package chapter3

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"runtime"
	"sync"
	"testing"
	"text/tabwriter"
	"time"
)

/*
	gorountine
	协程是非抢占式的、简单并发子 gorountine，不能被中断
*/

// gorountineExample gorountine demo
func gorountineExample() {
	sayHello := func() {
		fmt.Println("Hello world!")
	}

	go sayHello()
}

/*
	Go 语言的主机托管机制是一个名为 M:N 的调度器的实现，
M 个绿色线程映射到 N 个 OS 线程。然后 gorountine 安排在
绿色线程上。当 gorountine 数量超过可用的绿色线程时,调度
程序将处理分布在可用线程上的 gorountine，确保当这些 gorountine
被阻塞时，其他 gorountine 可以运行。
*/

/*
	Go 语言遵循一个成为 fork-join 的并发模型。
fork 指的是在程序中任意一点，它可以将执行的子分支与其父节点同时运行。
join 指的是在将来某个时候，这些并发的执行分支将会合并在一起。
*/

// syncExample sync 包 demo
func syncExample() {
	var wg sync.WaitGroup
	sayHello := func() {
		defer wg.Done()
		fmt.Println("Hello")
	}
	wg.Add(1)
	go sayHello()
	wg.Wait()
}

// syncExample2
// 闭包可以从创建它们的作用域中获取变量(的引用)
// 因此是结果是 "welcome "
func syncExample2() {
	var wg sync.WaitGroup
	salutation := "hello"
	wg.Add(1)
	go func() {
		defer wg.Done()
		salutation = "welcome"
	}()
	wg.Wait()
	fmt.Println(salutation)
}

// syncExample3
// 在任何 gorountine 开始运行之前，循环就会退出。
// 所以 salutation 会被转移到堆上，引用最后一个值 "good day"
// 所以输出结果：
//  good day
//  good day
//  good day
func syncExample3() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
	}
	wg.Wait()
}

// syncExample4
// 上述的函数的改正，将变量 salutation 的副本传递到包中，
// 输出结果：
//  hello
//  greetings
//  good day
func syncExample4() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait()
}

/*
	gorountine 非常轻量
	GC 并不会回收被丢弃的 gorountine
*/

// gorountineExample2 展示了了 gorountine 的内存占用
// 可能会影响性能的是上下文切换，即：当一个被托管的并发进程必须保存
// 它的状态以切换到一个不同的运行并发进程。如果并发进程太多，可能会
// 将所有 CPU 时间消耗在它们之间的上下文切换上，没有资源完成任何真正
// 需要 CPU 的工作。
func gorountineExample2() {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan interface{}
	var wg sync.WaitGroup
	noop := func() {
		wg.Done()
		<-c
	}

	const numGrountines = 1e4
	wg.Add(numGrountines)
	before := memConsumed()
	for i := numGrountines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed()
	fmt.Printf("%.3fkb", float64(after-before)/numGrountines/1000)
}

// contextSwitch 上下文切换展示
func contextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	var token struct{}
	sender := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			c <- token
		}
	}

	receiver := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			<-c
		}
	}

	wg.Add(2)
	go sender()
	go receiver()
	b.StartTimer()
	close(begin) // 两个 gorountine 开始运行
	wg.Wait()
}

/*
	Sync Package
*/

// waitGroupExample
// waitGroup 等待 gorountine 完成
func waitGroupExample() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("1st gorountine sleeping...")
		time.Sleep(1 * time.Second)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("2nd gorountine sleeping...")
		time.Sleep(2 * time.Second)
	}()

	wg.Wait()
	fmt.Println("All gorountines complete.")
	// waitGroup 可以看作一个并发—安全的计数器
	// 调用通过传入的整数执行 add 方法增加计数器的增量
	// 并调用 Done 方法对计数器进行增减，Wait 阻塞，直到计数器为零。
}

// waitGroupExample2
// 输出结果:
//	Hello from 5!
//	Hello from 4!
//	Hello from 3!
//	Hello from 2!
//	Hello from 1!
func waitGroupExample2() {
	hello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Printf("Hello from %v!\n", id)
	}

	const numGreeters = 5
	var wg sync.WaitGroup
	wg.Add(numGreeters)
	for i := 0; i < numGreeters; i++ {
		go hello(&wg, i+1)
	}
	wg.Wait()
}

/*
	互斥锁与读写锁
*/

// mutexExample
// 通过互斥锁对临界区保护
func mutexExample() {
	var count int
	var lock sync.Mutex

	increment := func() {
		lock.Lock()
		defer lock.Unlock()
		count++
		fmt.Printf("Incrementing: %d\n", count)
	}

	decrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count--
		fmt.Printf("Decrementing: %d\n", count)
	}

	// 增量
	var arithmetic sync.WaitGroup
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			increment()
		}()
	}

	// 减量
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			decrement()
		}()
	}

	arithmetic.Wait()
	fmt.Println("Arithmetic complete.")
}

// mutexExample2
func mutexExample2() {
	producer := func(wg *sync.WaitGroup, l sync.Locker) {
		// 第二个参数是 sync.Locker 类型，
		// 这个接口有两个方法 Lock 和 Unlock，
		// 分别对应 Mutex 和 RWMutex
		defer wg.Done()
		for i := 5; i > 0; i-- {
			l.Lock()
			l.Unlock()
			time.Sleep(time.Second)
		}
	}

	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}

	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1)
		beginTestTime := time.Now()
		go producer(&wg, mutex)
		for i := count; i > 0; i-- {
			go observer(&wg, rwMutex)
		}
		wg.Wait()
		return time.Since(beginTestTime)
	}

	var b byte
	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, b, 0)
	defer tw.Flush()

	var m sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutex\tMutex\n")
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			test(count, &m, m.RLocker()),
			test(count, &m, &m),
		)
	}
}

// condExample
// Cond 类型，一个 gorountine 的集合点，等待或发布一个 event
// 一个 "event" 是两个或两个以上的 gorountine 之间的任意信号
// c.Signal() 方法，它提供通知 gorountine 阻塞的调用 Wait，条件已经被触发。
// Signal 发现等待最长时间的 gorountine 并通知它。
// 另一个 boardcast() 方法，是向所有等待的 gorountine 发送信号。它提供了一种
// 同时与多个 gorountine 通信的方法。
// 与 channel 相比，Cond 类型的性能要高很多。
func condExample() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Removed from queue")
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(time.Second)
		c.L.Unlock()
	}
}

func condExample2() {
	type Button struct {
		Clicked *sync.Cond
	}

	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}
	subscribe := func(c *sync.Cond, fn func()) {
		// 允许我们注册函数处理来自条件的信号，每个处理程序都在自己的 gorountine 上运行
		// 并且订阅不会退出，直到 gorountine 被确认运行为止。
		var gorountineRunning sync.WaitGroup
		gorountineRunning.Add(1)
		go func() {
			gorountineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		gorountineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)
	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})

	button.Clicked.Broadcast()
	// 在 Clicked Cond 调用 Broadcast，所以三个处理程序都将运行

	clickRegistered.Wait()
}

/*
	Once 保证函数只调用一次
*/
func onceExample() {
	var count int
	increment := func() {
		count++
	}

	var once sync.Once
	var increments sync.WaitGroup
	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}

	increments.Wait()
	fmt.Printf("Count is %d\n", count)
}

/*
	Pool 池
*/

func poolExample() {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new instance.")
			return struct{}{}
		},
	}

	myPool.Get()
	instance := myPool.Get()
	myPool.Put(instance)
	myPool.Get()
}

// poolExample2 用 pool 可以节省内存
func poolExample2() {
	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}

	// 用 4kb 初始化 pool
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorks = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorks)
	for i := numWorks; i > 0; i-- {
		go func() {
			defer wg.Done()
			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}

// 用 pool 可以尽可能快地将预先分配的对象缓存加载启动
func connectToService() interface{} {
	time.Sleep(time.Second)
	return struct{}{}
}

func startNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()
		wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			connectToService()
			fmt.Fprintln(conn, "")
			conn.Close()
		}
	}()
	return &wg
}

func init() {
	daemonStarted := startNetworkDaemon()
	daemonStarted.Wait()
}

func BenchmarkNetworkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			b.Fatalf("cannot dial host: %v", err)
		}
		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}
		conn.Close()
	}
}

func warmServiceConnCache() *sync.Pool {
	p := &sync.Pool{
		New: connectToService,
	}

	for i := 0; i < 10; i++ {
		p.Put(p.New())
	}

	return p
}

func startNetworkCacheDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		connPool := warmServiceConnCache()
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()
		wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
			}
			svcConn := connPool.Get()
			fmt.Fprintln(conn, "")
			connPool.Put(svcConn)
			conn.Close()
		}
	}()
	return &wg
}

/*
	channel, 充当信息的传送管道，值可以沿着 channel 传递
*/

// chanExample channel 用法举例
func chanExample() {
	var dataStream chan interface{}     // 声明
	dataStream = make(chan interface{}) // 实例化
	_ = dataStream
}

// channel 可以只读取，也可以只发送
func chanExample2() {
	var dataStream <-chan interface{} // 只能读取的 channel
	dataStream = make(<-chan interface{})
	_ = dataStream

	var dataStream2 chan<- interface{}
	dataStream2 = make(chan<- interface{})
	_ = dataStream2
}

// chanExample3
// Go 语言会隐式地将双向 channel 转换为单向channel
func chanExample3() {
	var receiveChan <-chan interface{}
	var sendChan chan<- interface{}
	dataStream := make(chan interface{})

	// 有效的语法
	receiveChan = dataStream
	sendChan = dataStream

	_ = receiveChan
	_ = sendChan
}

// chanExample4
func chanExample4() {
	stringStream := make(chan string)
	go func() {
		stringStream <- "Hello channels" // 将字符串文本传递给 stringStream
	}()
	fmt.Println(<-stringStream) // 读取channel里的字符串并打印
}

// 将一个值写入只读的 channel 是错误的 -> 引发异常(invalid operation)
// 从一个只可以写的 channel 读取值也是错误的 -> 引发异常(invalid operation)
// Go 语言中 channel 式阻塞的。只有 channel 中的数据被消费后，新的数据才能被写入。
// 任何试图从空 channel 读取数据的 gorountine 将等待至少一条数据被写入 channel 后
// 才能读到。 fmt.Println 会从 stringStream 这个 channel 中消费一条数据。所以，
// 它会等 channel 中有数据后才开始消费。

// 同样，匿名 gorountine 试图向 stringStream 里写入一条字符串，所以在写入成功之前
// gorountine 将不会退出。因此 main gorountine 和匿名 gorountine 一定是阻塞住的。

func chanExample5() {
	stringStream := make(chan string)

	go func() {
		if 0 != 1 { // 一定会触发条件，stringStream 不会被写入，引发 panic 死锁
			return
		}
		stringStream <- "Hello channels!"
	}()
	fmt.Println(<-stringStream)
}

// chanExample6
// <- 操作符可以返回两个值
// 输出结果：
//	(true): Hello channels!
func chanExample6() {
	stringStream := make(chan string)
	go func() {
		stringStream <- "Hello channels!"
	}()
	salutation, ok := <-stringStream
	fmt.Printf("(%+v): %v", ok, salutation)
}

// ok 表示该 channel 上有新数据写入，
// 或者是有 closed channel 生成的默认值
// 能够提示 channel 中是否会有新的值写入，有助于下游程序
// 知道什么时候消费、退出、给新的 channel 重新建立连接等。
// channel 可以使用 close 关闭

// chanExample7
// 从一个已经关闭的 channel 读取数据
// 输出结果：
//	(false): 0
func chanExample7() {
	intStream := make(chan int)
	close(intStream)
	integer, ok := <-intStream
	fmt.Printf("(%v): %v", ok, integer)
}

// channel 没有写入数据，直接被 close，仍旧能够被读取
// 这是为了支持一个 channel 有单个上游写入，有多个
// 下游读取。

// chanExample8
//
func chanExample8() {
	intStream := make(chan int)
	go func() {
		defer close(intStream) // 在 gorountine 退出之前，确保 channel 是关闭的
		for i := 1; i <= 5; i++ {
			intStream <- i
		}
	}()

	for integer := range intStream { // 使用 range 遍历
		fmt.Printf("%v", integer)
	}
}

// 这里的循环不需要退出条件，并且 range 方法不会返回第二个布尔值。
// 处理一个已关闭的 channel 的细节可以让你保持循环清洁。

// chanExample9
// 关闭 channel 也是一种同时给多个 gorountine 发信号的方法。
// 如果有 n 个 gorountine 在一个 channel 上等待，
// 而不是在 channel 上写 n 次来打开每个 goroutine，你可以简单地关闭 channel。
// 由于一个被关闭的 channel 可以被无数次读取，关闭比执行 n 次更适合，也更快。
func chanExample9() {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin // goroutine 一直会在这里等待，直到 begin 不再阻塞
			fmt.Printf("%v has begun\n", i)
		}(i)
	}
	fmt.Println("Unblocking gorountines...")
	close(begin) // 关闭 channel，从而同时打开所有阻塞着的 gorountine
	wg.Wait()
}

// buffered channel，提供容量的 channel
// 意味着即使没有在 channel 上执行读取操作，gorountine 仍然可以执行 n 写入
// n 是缓冲 channel 的容量
func chanExample10() {
	var dataStream chan interface{}
	dataStream = make(chan interface{}, 4)
	_ = dataStream
}

// 当讨论阻塞时，如果说 channel 是满的，那么 channel 阻塞。
// 缓冲 channel 是一个内存中的 FIFO 队列，用于并发进程进行通信。

// bufferedChanExample 缓冲 channel 的示例
// 输出结果：
//	Sending: 0
//	Sending: 1
//	Sending: 2
//	Sending: 3
//	Sending: 4
//  Sending Done.
//  Received 0.
//  Received 1.
//  Received 2.
//  Received 3.
//  Received 4.
func bufferedChanExample() {
	var stdouBuff bytes.Buffer         // 内存缓冲区，比直接写 stdout 快
	defer stdouBuff.WriteTo(os.Stdout) // 程序退出前，缓冲区内容写入 stdout

	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdouBuff, "Producer Done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdouBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Fprintf(&stdouBuff, "Received %v.\n", integer)
	}
}

// channel 各种状态在读/写下的结果

//         Channel状态   结果
//
// Read    nil           阻塞
//         打开但非空    输出值
//         打开但空      阻塞
//         关闭的        <默认值>, false
//         只写          编译错误

// Write   nil           阻塞
//         打开但填满    阻塞
//         打开但不满    写入值
//         关闭的        panic
//         只读          编译错误

// Close   打开但非空    关闭Channel：读取成功，直到通道耗尽，然后读取产生值的默认值
//         打开但空      关闭Channel：读到生产者的默认值
//		   关闭的        panic
//         只读          编译错误
