package chapter3

import (
	"io/ioutil"
	"net"
	"testing"
)

func init() {
	daemonStarted := startNetworkDaemon()
	daemonStarted.Wait()

	daemonStartedConnCache := startNetworkDaemonConnCache()
	daemonStartedConnCache.Wait()
}

// goroutineExample : goroutine 并发 demo
// go test ./chapter3 -run TestGoroutineExample -v
func TestGoroutineExample(t *testing.T) {
	goroutineExample()
}

// waitGroupExample: WaitGroup 的用法 demo
// go test ./chapter3 -run TestWaitGroupExample -v
func TestWaitGroupExample(t *testing.T) {
	waitGroupExample()
}

// mutextExample 互斥锁用法 demo
// go test ./chapter3 -run mutextExample -v
func TestMutextExample(t *testing.T) {
	mutextExample()
}

// syncCondExample sync.Cond 用法 demo
// go test ./chapter3 -run syncCondExample -v
func TestSyncCondExample(t *testing.T) {
	syncCondExample()
}

// boardcastExample sync cond boardcast 用法 demo
// go test ./chapter3 -run boardcastExample -v
func TestBoardcaskExample(t *testing.T) {
	boardcastExample()
}

// syncOnceExample sync.Once 用法 demo
// go test ./chapter3 -run syncOnceExample -v
func TestSyncOnceExample(t *testing.T) {
	syncCondExample()
}

// go test ./chapter3 -benchtime=10s -bench=. -run=none -v
func BenchmarkNetworkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			b.Fatal(err)
			return
		}
		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %+v", err)
		}
		conn.Close()
	}
}

func BenchmarkNetworkRequestConnCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8081")
		if err != nil {
			b.Fatal(err)
		}
		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %+v", err)
		}
		conn.Close()
	}
}
