package channel

import "testing"

// go test ./channel -run TestPiplineExample -v -count=1
func TestPiplineExample(t *testing.T) {
	piplineExample()
}

// go test ./channel -run TestPiplineExample2 -v -count=1
func TestPiplineExample2(t *testing.T) {
	piplineExample2()
}

// go test ./channel -run TestPiplineExample3 -v -count=1
func TestPiplineExample3(t *testing.T) {
	piplineExample3()
}

// go test ./channel -run TestPiplineExample4 -v -count=1
func TestPiplineExample4(t *testing.T) {
	piplineExample4()
}
