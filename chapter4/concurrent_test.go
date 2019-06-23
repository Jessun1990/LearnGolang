package chapter4

import "testing"

/*
go test ./chapter4 -run TestConCurrentExample -v
*/
func TestConCurrentExample(t *testing.T) {
	concurrentExmaple()
}

/*
go test ./chapter4 -run TestConCurrentExample2
*/
func TestConCurrentExample2(t *testing.T) {
	concurrentExmaple2()
}
