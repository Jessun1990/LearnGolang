// Package chapter2 对你的代码建模：通信顺序进程
// 第二章主要讲解Go在并发上的处理原则
package chapter2

// 并发和并行
// 并发属于代码，并行属于一个运行中的程序

// Go 并发： CSP 模型 --> goroutine， 也支持“锁的方式来并发
// 不要通过共享内存来通信，而应该通过通信来共享内存
// 意思就是：多用 channel ，少用 sync mutex
