package chapter3

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/hsyan2008/go-logger"
)

// goroutineExample1 : goroutine 并发 demo
// go test ./chapter3 -run TestGoroutineExample -v
func goroutineExample() {
	for i := 0; i <= 10; i++ {
		go func(no int) {
			fmt.Printf("i = %+v\n", no)
		}(i)
		time.Sleep(time.Millisecond)
	}
}

// waitGroupExample1: WaitGroup 的用法 demo
// go test ./chapter3 -run TestWaitGroupExample -v
func waitGroupExample() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("1st goroutine sleeping...")
		time.Sleep(1 * time.Second)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("2nd goroutine sleeping...")
		time.Sleep(2 * time.Second)
	}()
	wg.Wait()
	fmt.Println("All goroutine complete.")
}

// mutextExample 互斥锁用法 demo
// go test ./chapter3 -run mutextExample -v
func mutextExample() {
	var count int
	var lock sync.Mutex

	increment := func() {
		lock.Lock()
		defer lock.Unlock()
		count++
		fmt.Printf("Incrementing: %+v\n", count)
	}

	decrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count--
		fmt.Printf("Decrementing: %+v\n", count)
	}

	var arithmetic sync.WaitGroup
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			increment()
		}()
	}

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

// syncCondExample sync.Cond 用法 demo
// go test ./chapter3 -run syncCondExample -v
func syncCondExample() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10) // 最终会添加10个项目，所以用10的容量实例化

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Removed from queue")
		c.L.Lock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
		}

		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(time.Second) // 创建一个新 goroutine，在一秒钟后删除一个元素
		c.L.Unlock()
	}
}

// boardcastExample sync cond boardcast 用法 demo
// go test ./chapter3 -run boardcastExample -v
func boardcastExample() {
	type Button struct { // 模拟一个 Button 类型
		Clicked *sync.Cond
	}
	button := Button{
		Clicked: sync.NewCond(&sync.Mutex{}),
	}
	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)
	subscribe(button.Clicked, func() {
		fmt.Println("Maxmizing window.")
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

	button.Clicked.Broadcast() // 使用 Broadcast，所有三个处理程序都将运行
	clickRegistered.Wait()
}

// syncOnceExample sync.Once 用法 demo
// go test ./chapter3 -run syncOnceExample -v
func syncOnceExample() {
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

// syncPoolExample sync pool 用法 demo
// go test ./chapter3 -run syncPoolExample -v
func syncPoolExample() {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new instance")
			return struct{}{}
		},
	}

	myPool.Get()
	instance := myPool.Get()
	myPool.Put(instance)
	myPool.Get()
}

// startNetworkDaemon 模拟对服务端的请求连接
func startNetworkDaemon() *sync.WaitGroup {
	connectToService := func() interface{} {
		return struct{}{}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			logger.Warn(err)
			return
		}
		defer server.Close()
		for {
			conn, err := server.Accept()
			if err != nil {
				logger.Warn(err)
				continue
			}
			connectToService()
			fmt.Fprintln(conn, "")
			conn.Close()
		}
	}()
	return &wg
}

// startNetworkDaemonConnCache 模拟对服务端的请求连接
func startNetworkDaemonConnCache() *sync.WaitGroup {
	connectToService := func() interface{} {
		return struct{}{}
	}

	warmServiceConnCache := func() *sync.Pool {
		p := &sync.Pool{
			New: connectToService,
		}

		for i := 0; i < 10; i++ {
			p.Put(p.New())
		}
		return p
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		connPool := warmServiceConnCache()
		server, err := net.Listen("tcp", "localhost:8081")
		if err != nil {
			logger.Warn(err)
			return
		}
		defer server.Close()
		wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				logger.Warn(err)
				return
			}
			svcConn := connPool.Get()
			fmt.Fscanln(conn, "")
			connPool.Put(svcConn)
			conn.Close()
		}
	}()
	return &wg
}
