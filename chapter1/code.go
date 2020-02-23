// Package chapter1 并发概述
// 本书第一章主要讲解并发的难点以及特点

package chapter1

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*
	并发的难点和特点
*/
// raceCondCase：竞争条件的代码举例
func raceCondCase() {
	var data int
	go func() {
		data++
	}()
	if data == 0 {
		fmt.Printf("the value is %v. \n", data)
	}
}

// atomicity 原子性

// memory Access Synchronization：内存访问同步
func memAccessSyncCase() {
	var data int
	go func() {
		data++
	}()
	if data == 0 {
		fmt.Println("the value is 0.")
	} else {
		fmt.Printf("the value is %v.", data)
	}
}

// memAccessSyncMtxCase： 内存访问同步之锁的应用
func memAccessSyncMtxCase() {
	var memAccess sync.Mutex
	var value int
	go func() {
		memAccess.Lock()
		value++
		memAccess.Unlock()
	}()
	memAccess.Lock()
	if value == 0 {
		fmt.Printf("the value is %v. \n", value)
	} else {
		fmt.Printf("the value is %v. \n", value)
	}
	memAccess.Unlock()
}

// deadLockCase：死锁
// 死锁程序是所有并发进程彼此等待的程序。在这种情况下
// 没有外界的干预，这个程序将永远无法恢复。
func deadLockCase() {
	var wg sync.WaitGroup
	printSum := func(v1, v2 *value) {
		defer wg.Done()
		v1.mtx.Lock()
		defer v1.mtx.Unlock()

		time.Sleep(2 * time.Second)
		v2.mtx.Lock()
		defer v2.mtx.Unlock()

		fmt.Printf("sum=%v\n", v1.value+v2.value)
	}

	var a, b value
	wg.Add(2)
	go printSum(&a, &b)
	go printSum(&b, &a)
	/*
	   第一次调用 printSum 锁定 a，然后试图锁定 b
	   第二次调用 printSum 已锁定 b，并试图锁定 a
	   这两个 gorountine 都无限等待
	*/
	wg.Wait()
}

type value struct {
	mtx   sync.Mutex
	value int
}

/*

    构成死锁的 Conffman 条件：
1. 相互排斥：并发进程同时拥有资源的独占权。
2. 等待条件：并发进程必须同时拥有一个资源，并等待额外的资源。
3. 没有抢占：并发进程拥有的资源只能被该进程释放，即可满足这个条件。
4. 循环等待：一个并发进程P1 必须等待一系列其他并发进程P2，
   这些并发进程，同时也在等待进程P1，这样便满足了这个最终条件。


*/

// liveLockCase：活锁
// 活锁是正在主动执行并发操作的程序，但是这些操作无法向前推荐程序的状态。
func liveLockCase() {
	cadence := sync.NewCond(&sync.Mutex{})

	go func() {
		for range time.Tick(1 * time.Microsecond) {
			cadence.Broadcast() // sync.Cond{}.Broadcast() 函数，
		}
	}()

	takeStep := func() {
		cadence.L.Lock()
		cadence.Wait() // sync.Cond{}.Wait() 函数，所有 gorountine 等待通知
		cadence.L.Unlock()
	}

	// tryDir 允许一个人尝试向一个方向移动，并返回是否成功。
	// Dir，每个方向都表示为试图朝向这个方向移动的人数
	tryDir := func(dirName string, dir *int32, out *bytes.Buffer) bool {
		fmt.Fprintf(out, "%v", dirName)
		atomic.AddInt32(dir, 1) // 原子级操作，1 赋值给 dir
		takeStep()
		if atomic.LoadInt32(dir) == 1 {
			fmt.Fprint(out, ". Success!")
			return true
		}

		takeStep()
		atomic.AddInt32(dir, -1)
		return false
	}

	var left, right int32
	tryLeft := func(out *bytes.Buffer) bool { return tryDir("left", &left, out) }
	tryRight := func(out *bytes.Buffer) bool { return tryDir("right", &right, out) }

	_ = tryLeft
	_ = tryRight

	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer
		defer func() { fmt.Println(out.String()) }()
		defer walking.Done()

		fmt.Fprintf(&out, "%v is trying to scoot: ", name)
		for i := 0; i < 5; i++ {
			if tryLeft(&out) || tryRight(&out) {
				return
			}
		}
		fmt.Fprintf(&out, "\n%v tosses her hands up in exaspertion!", name)
	}

	var peopleInHallway sync.WaitGroup
	peopleInHallway.Add(2)
	go walk(&peopleInHallway, "Alice")
	go walk(&peopleInHallway, "Barbara")
	peopleInHallway.Wait()
}

// starvationCase 饥饿
// 饥饿是在任何情况下，并发进程都无法获得执行工作所需的所有资源
func starvationCase() {
	var wg sync.WaitGroup
	var shareLock sync.Mutex
	const runtime = 5 * time.Second

	greedWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) < runtime; {
			shareLock.Lock()
			time.Sleep(3 * time.Nanosecond)
			shareLock.Unlock()
			count++
		}
		fmt.Printf("Greedy worker was able to execute %v wokr loops\n", count)
	}

	politeWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			shareLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			shareLock.Unlock()

			shareLock.Lock()
			time.Sleep(time.Nanosecond)
			shareLock.Unlock()

			shareLock.Lock()
			time.Sleep(time.Nanosecond)
			shareLock.Unlock()

			count++
		}
		fmt.Printf("Polite worker was able to execute %v work loops\n", count)
	}

	wg.Add(2)
	go greedWorker()
	go politeWorker()

	wg.Wait()
}
