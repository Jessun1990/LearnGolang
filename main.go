package main

import (
	//"github.com/jessun1990/LearnGolang/wg"
	//"fmt"
	"github.com/jessun1990/LearnGolang/channel"
)

func main() {
	//wg.TryWg()
	//stringStream := make(chan string)
	//go func() {
	//stringStream <- "hello channels"
	//}()
	//fmt.Print(<-stringStream)

	channel.TryGoroutine1()
}
