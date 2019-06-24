package chapter4

import "testing"

/*
go test ./chapter4 -run TestErrHandleExample -v
*/
func TestErrHandleExample(t *testing.T) {
	errHandleExample()
}

/*
go test ./chapter4 -run TestErrHandleExample2 -v
*/
func TestErrHandleExample2(t *testing.T) {
	errHandleExample2()
}
