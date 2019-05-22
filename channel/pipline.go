package channel

import "fmt"

// Pipline 的组合用法
func piplineExample() {
	multiply := func(values []int, multiplier int) []int {
		multipiedValues := make([]int, len(values))
		for i, v := range values {
			multipiedValues[i] = v * multiplier
		}
		return multipiedValues
	}

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}

	ints := []int{1, 2, 3, 4}
	for _, v := range multiply(add(multiply(ints, 2), 1), 2) {
		fmt.Println(v)
	}

}

func piplineExample2() {
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
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

	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int,
	) <-chan int {
		multipiedStream := make(chan int)
		go func() {
			defer close(multipiedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipiedStream <- i * multiplier:

				}
			}
		}()
		return intStream
	}

	add := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
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

	// Pipline 最佳实践
	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)
	pipline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipline {
		fmt.Println(v)
	}
}

// 生成器举例
func piplineExample3() {

	repeat := func(done <-chan interface{}, values ...interface{},
	) <-chan interface{} {
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

	//done := make(chan interface{})
	//close(done)

	//for x := range repeat(done, []int{1, 2, 3}) {
	//fmt.Println(x)
	//}

	take := func(done <-chan interface{}, valueStream <-chan interface{},
		num int,
	) <-chan interface{} {

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
		fmt.Printf("%+v\n", num)
	}

}
