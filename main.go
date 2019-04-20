package main

import (
	sync "github.com/jessun1990/LearnGolang/sync"
)

func main() {
	//sayHello := func() {
	//fmt.Println("hello")
	//}

	//go sayHello()
	/*
	   输出没有能在在返回前完成，所有显示为空。
	*/

	sync.TryRWMutex()
}
