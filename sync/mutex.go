package sync

// golang 中的 sync 实现了两种锁:
// Mutex: 互斥锁
// RWMutex: 读写锁，RWMutex 基于 Mutex 实现

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

/*
- TryMutex 互斥锁举例，保证对资源的独占访问
- Mutex 为互斥锁，Lock() 为加锁，Unlock() 为解锁
- 在一个 goroutine 获得 Mutex 后，其他 goroutine 只能等到这个 goroutine 释放该 Mutex
- 在使用 Lock() 加锁后，不能在继续加锁，直到利用 Unlock() 解锁后才能加锁
- 在 Lock() 之前使用 Unlock() 会导致 panic 异常
- 已经锁定的 Mutex 并不与特定的 goroutine 相关联，这样可以利用一个 goroutine 对其加锁，
  再利用其他的 goroutine 对其解锁
- 在同一个 goroutine 中的 Mutex 解锁之前再次对其进行加锁，会导致死锁
- 适用于读写不确定，并且只有一个读或者写的场景
*/
func TryMutex() {
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
}

/*
- TryRWMutex 读写锁
- RWMutex 是单写多读锁，该锁可以加多个读锁或者一个写锁
- 读锁占用的情况下会阻止写，不会阻止读，多个 goroutine 可以同时获取读锁
- 写锁会阻止其他 goroutine （无论读写） 进来，整个锁由该 goroutine 独占
- 适用于读多写少的的场景

- Lock() 加写锁，Unlock() 解写锁
- 如果在加写锁之前已经有其他的读锁和写锁，则 Lock() 会阻塞直到该锁可用，
  为确保该锁可用，已经阻塞的 Lock() 调用会从获得的锁中排除新的读取器，
  即写锁权限高于读锁，有写锁时优先进行写锁定
- 在 Lock() 之前使用 Unlock() 会导致 panic 异常

- RLock() 和 RUnlock()
- RLock() 加读锁时，如果存在写锁，则无法接续加锁；当只有读锁或者没有锁时，
  可以加读锁，读锁可以加载多个
- RUnlock() 解读锁，RUnlock() 撤销单次 RLock() 调用，对于其他同时存在的读
  锁则没有效果
- 在没有读锁的情况下调用 RUnlock() 会导致 panic 错误
- RUnlock() 的个数不得多余 RLock()，否则会导致 panic 错误
*/
func TryRWMutex() {
	producer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		for i := 5; i > 0; i-- {
			l.Lock()
			l.Unlock()
			time.Sleep(1)
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

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
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
