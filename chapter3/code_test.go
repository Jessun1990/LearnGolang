package chapter3

import "testing"

/*
go test ./chapter3 -v -count=1 -run TestChanExmaple9
*/
func TestChanExmaple9(t *testing.T) {
	chanExample9()
}

/*
go test ./chapter3 -v -count=1 -run TestSelectExample2
*/
func TestSelectExample2(t *testing.T) {
	selectExample2()
}

/*
go test ./chapter3 -v -count=1 -run TestSelectExample3
*/
func TestSelectExample3(t *testing.T) {
	selectExample3()
}

/*
go test ./chapter3 -v -count=1 -run TestSelectExample6
*/
func TestSelectExample6(t *testing.T) {
	selectExample6()
}
