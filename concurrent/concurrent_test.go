package concurrent

import "testing"

// go test ./concurrent/ -run TestGoroutineExample1 -v
func TestGoroutineExample(t *testing.T) {
	goroutineExample()
}

// go test ./concurrent/ -run TestGoroutineExample2 -v
func TestGoroutineExample2(t *testing.T) {
	goroutineExample2()
}

// go test ./concurrent/ -run TestGoroutineExample3 -v
func TestGoroutineExample3(t *testing.T) {
	goroutineExample3()
}

// go test ./concurrent/ -run TestGoroutineExample4 -v
func TestGoroutineExample4(t *testing.T) {
	goroutineExample4()
}
