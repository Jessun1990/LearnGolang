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

// go test ./concurrent/ -run TestGoroutineExample5 -v
func TestGoroutineExample5(t *testing.T) {
	goroutineExample5()
}

// go test ./concurrent/ -run TestGoroutineExample6 -v
func TestGoroutineExample6(t *testing.T) {
	goroutineExample6()
}
