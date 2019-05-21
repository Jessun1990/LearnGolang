package channel

import "testing"

// go test ./channel -run TestChanExample1 -v
func TestChanExample1(t *testing.T) {
	chanExample1()
}

// go test ./channel -run TestChanExample2 -v
func TestChanExample2(t *testing.T) {
	chanExample2()
}

// go test ./channel -run TestChanExample3 -v
func TestChanExample3(t *testing.T) {
	chanExample3()
}

// go test ./channel -run TestChanExample4 -v
func TestChanExample4(t *testing.T) {
	chanExample4()
}

// go test ./channel -run TestChanExample5 -v
func TestChanExample5(t *testing.T) {
	chanExample5()
}

// go test ./channel -run TestChanExample6 -v
func TestChanExample6(t *testing.T) {
	selectExample1()
}

// go test ./channel -run TestChanExample7 -v
func TestChanExample7(t *testing.T) {
	selectExample2()
}

func TestChanExample8(t *testing.T) {
	selectExample3()
}

func TestChanExample9(t *testing.T) {
	selectExample4()
}
