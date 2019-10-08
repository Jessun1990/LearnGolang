package chapter3

//import (
//"fmt"
//"time"
//)

//// selectExample : select 用法 demo
//func selectExample() {
//start := time.Now()
//c := make(chan struct{})
//go func() {
//time.Sleep(5 * time.Second)
//close(c)
//}()
//fmt.Println("Blocking on read...")
//select {
//case <-c:
//fmt.Printf("Unblocking %+v laster.", time.Since(start))
//}
//}

//// selectExample2 : 多个 channel 同时可用，就随机选取执行
//func selectExample2() {
//c1 := make(chan interface{})
//close(c1)
//c2 := make(chan interface{})
//close(c2)

//var c1Count, c2Count int
//for i := 1000; i >= 0; i-- {
//select {
//case <-c1:
//c1Count++
//case <-c2:
//c2Count++
//}
//}
//fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
//}

//// selectExample3 ：select 中 default 用法 demo
//func selectExample3() {
//start := time.Now()
//var c <-chan int
//select {
//case <-c:
//case <-time.After(time.Second):
//fmt.Println("Time out.")
//default:
//fmt.Printf("In default after %+v\n", time.Since(start))
//}
//}

//// selectExample4 ： 加入 loop 标签来配合 channel 控制流程
//func selectExample4() {
//done := make(chan interface{})
//go func() {
//time.Sleep(5 * time.Second)
//close(done)
//}()

//workCounter := 0
//loop:
//for {
//select {
//case <-done:
//break loop
//default:
//fmt.Println("This is default")
//}
//workCounter++
//time.Sleep(time.Second)
//}
//fmt.Printf("Achieved %+v cycles of work before signalled to stop. \n", workCounter)
//}

//// selectExample5：将任意数量的 channel 组合到单个 channel 中，只要任何
//// 子 channel 关闭或者写入，该 channel 就会关闭
//func selectExample5() {
//var or func(chans ...<-chan interface{}) <-chan interface{}

//or = func(chans ...<-chan interface{}) <-chan interface{} {
//switch len(chans) {
//case 0:
//return nil
//case 1:
//return chans[0]
//}

//orDone := make(chan interface{})

//go func() {
//defer close(orDone)
//switch len(chans) {
//case 2:
//select {
//case <-chans[0]:
//case <-chans[1]:
//}
//default:
//select {
//case <-chans[0]:
//case <-chans[1]:
//case <-chans[2]:
//case <-or(append(chans[3:], orDone)...):
//}
//}
//}()
//return orDone
//}

//sig := func(after time.Duration) <-chan interface{} {
//c := make(chan interface{})
//go func() {
//defer close(c)
//time.Sleep(after)
//}()
//return c
//}

//start := time.Now()
//<-or(
//sig(2*time.Hour),
//sig(2*time.Second),
//sig(2*time.Minute),
//)
//fmt.Printf("Done after %+v", time.Since(start))
//}

//// chanExample6 ：使用 channel 并发时的错误处理
////func chanExample6() {
////type customRes struct {
////Error    error
////Response *http.Response
////}

////checkStatus := func(done <-chan interface{},
////urls ...string) <-chan customRes {
////results := make(chan customRes)

////go func() {
////defer close(results)
////for _, url := range urls {
////rsp, err := http.Get(url)
////res := customRes{
////Error:    err,
////Response: rsp,
////}
////select {
////case <-done:
////return
////case results <- res:
////}
////}
////}()
////return results
////}

////done := make(chan interface{})
////defer close(done)

////urls := []string{"https://www.google.com", "https://badhost"}
////for result := range checkStatus(done, urls...) {
////if result.Error != nil {
////fmt.Printf("error: <%+v>", result.Error)
////continue
////}
////fmt.Printf("Response: %+v\n", result.Response.Status)
////}
////}
