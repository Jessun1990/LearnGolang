/*
 sync.Once 只调用一次，不再考虑其他并发情况。
*/
package sync

import (
	"fmt"
	"sync"
)

func TryOnce() {
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
	fmt.Println("Count is %d\n", count)
}
