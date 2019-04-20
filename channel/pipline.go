package pipline

import (
	"fmt"
)

// TryPipline1 ...
// 错误用例
func TryPipline1() {
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

	ints := []int{1, 2, 3, 4}
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}
}

func TryPipline2() {
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intSteam := make(chan int)
		go func() {
			defer close(intSteam)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intSteam <- i:
				}
			}
		}()
		return intSteam
	}

	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedSteam := make(chan int)
		go func() {
			defer close(multipliedSteam)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedSteam <- i * multiplier:
				}
			}
		}()
		return multipliedSteam
	}

	add := func(done <-chan interface{}, intSteam <-chan int, additive int) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intSteam {
				select {
				case <-done:
					return
				case addedStream <- i + additive:
				}
			}
		}()
		return addedStream
	}

	// main
	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)
	pipline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)
	for v := range pipline {
		fmt.Println(v)
	}
}
