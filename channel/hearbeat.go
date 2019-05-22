package channel

import (
	"fmt"
	"time"
)

func heartbeatExample() {
	doWork := func(done <-chan interface{}, pulseInterval time.Duration) (

		<-chan interface{}, <-chan time.Time) {

		heartbeat := make(chan interface{})
		// 设置一个发送心跳型号的通道。doWork 会返回该通道
		results := make(chan time.Time)

		go func() {
			defer close(heartbeat)
			defer close(results)

			pulse := time.Tick(pulseInterval)       // pulseInterval 心跳的间隔时间
			workGen := time.Tick(2 * pulseInterval) // 用来模拟工作

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default: // 没有接受到心跳的情况
				}
			}
			sendResult := func(r time.Time) {
				for {
					select {
					case <-done:
						return
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			for {
				select {
				case <-done:
					return
				case <-pulse: // 就像done channel 一样，当你执行发送或者接收时,也需要包含发送心跳的分支
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	}

	_ = doWork

	// ==== 消费
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })
	// 声明了一个标准的 done channel，并在10秒钟后关闭。

	const timeout = 2 * time.Second               // 超时时间
	heartbeat, results := doWork(done, timeout/2) // timeout/2 为心跳时间

	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")

		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %+v\n", r.Second())
		case <-time.After(timeout):
			return

		}
	}

}
