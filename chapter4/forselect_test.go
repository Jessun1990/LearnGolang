package chapter4

import "testing"

/*
go test ./chapter4 -run TestGoroutineExample2 -v
*/
func TestGoroutineExample2(t *testing.T) {
	goroutineExample2()
}

/*
go test ./chapter4 -run TestGoroutineExample3 -v
*/
func TestGoroutineExample3(t *testing.T) {
	goroutineExample3()
}
