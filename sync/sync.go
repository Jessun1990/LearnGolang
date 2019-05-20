package sync

import (
	"fmt"
	"net"
	"sync"
	"time"

	logger "github.com/hsyan2008/go-logger"
)

// ShowWaitGroup : go test ./sync  -run TestShowWaitGroup -v
func ShowWaitGroup() {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("1st goroutine sleeping ...")
		time.Sleep(1)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("2nd goroutine sleeping ...")
		time.Sleep(2)
	}()
	wg.Wait()
	fmt.Println("All goroutine complete.")
}

// ShowMutex :go test ./sync  -run TestShowMutex -v
func ShowMutex() {
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
	fmt.Println("Arithmetic complete")
}

// ShowSyncCond : go test ./sync -run TestShowSyncCond
func ShowSyncCond() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue :=
		func(delay time.Duration) {
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
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
	}
}

// ShowSyncOnce :
func ShowSyncOnce() {
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

// ****************************************************************
// ShowSyncPool :
func ShowSyncPool() {
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

// sync.Pool 性能变化
func connectToService() interface{} { // connectToService 模拟创建到服务的连接
	time.Sleep(time.Second)
	return struct{}{}
}

func startNetworkDaemon() *sync.WaitGroup { // 如果服务为每个请求都启动一个新的连接。
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		server, err := net.Listen("tcp", "localhost:8080")
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
				continue
			}
			connectToService()
			fmt.Fprintln(conn, "")
			conn.Close()
		}

	}()

	return &wg
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

func startNetworkDaemonConnCache() *sync.WaitGroup {
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
				continue
			}
			svcConn := connPool.Get()
			fmt.Fscanln(conn, "")
			connPool.Put(svcConn)
			conn.Close()
		}
	}()
	return &wg
}
