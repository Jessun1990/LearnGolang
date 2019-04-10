package channel

import (
	"fmt"
	"time"
)

func TrySelect1() {

	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()

	fmt.Println("Blocking on read...")

	select {
	case <-c:
		fmt.Printf("unblocked %v later.\n", time.Since(start))
	}
}

func TrySelect2() {
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

func TrySelect3() {
	var c <-chan int
	select {
	case <-c:
	case <-time.After(1. * time.Second):
		fmt.Println("Time out.")

	}
}
