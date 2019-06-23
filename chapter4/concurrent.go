package chapter4

import "fmt"

// 特定约束举例 P100
func concurrentExmaple() {
	data := make([]int, 4)
	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}
	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}

// 词法约束举例 P101
func concurrentExmaple2() {
	chanOwer := func() <-chan int {
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Printf("Received: %+v\n", result)
		}
		fmt.Println("Done receiving")
	}

	results := chanOwer()
	consumer(results)
}
