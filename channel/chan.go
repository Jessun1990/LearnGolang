package channel

import (
	"fmt"
)

func TryChan() {

	chanOwner := func() <-chan int {
		resultsStream := make(chan int, 5)
		go func() {
			defer close(resultsStream)
			for i := 0; i <= 5; i++ {
				resultsStream <- i
			}
		}()
		return resultsStream
	}

	resultsSteam := chanOwner()
	for r := range resultsSteam {
		fmt.Printf("Received: %d\n", r)
	}
	fmt.Println("Done receiving!")

}
