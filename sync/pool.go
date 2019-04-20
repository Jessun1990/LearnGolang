package sync

import (
	"fmt"
	"sync"
)

// Pool ...
func Pool() {
	myPool := sync.Pool{
		New: func() interface{} {
			fmt.Println("creating new instance")
			return struct{}{}
		},
	}
	myPool.Get()
}
