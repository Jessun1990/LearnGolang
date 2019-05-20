package sync

import (
	"io/ioutil"
	"net"
	"testing"
)

// go test ./sync  -run TestShowWaitGroup -v
func TestShowWaitGroup(t *testing.T) {
	ShowWaitGroup()
}

// go test ./sync  -run TestShowMutex -v
func TestShowMutex(t *testing.T) {
	ShowMutex()
}

// go test ./sync -run TestShowOnce -v
func TestShowOnce(t *testing.T) {
	ShowSyncOnce()
}

// go test ./sync -run TestSyncCond
func TestShowSyncCond(t *testing.T) {
	ShowSyncCond()
}

func TestShowSyncPool(t *testing.T) {
	ShowSyncPool()
}

func init() {
	daemonStarted := startNetworkDaemon()
	daemonStarted.Wait()

	daemonStartedConnCache := startNetworkDaemonConnCache()
	daemonStartedConnCache.Wait()
}

// go test ./sync -benchtime=10s -bench=. -run=none -v
func BenchmarkNetworkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			b.Fatal(err)
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

// output :
//goos: linux
//goarch: amd64
//pkg: concurrency_in_go/sync
//BenchmarkNetworkRequest-12                    10        1001102219 ns/op
//BenchmarkNetworkRequestConnCache-12        11581           1523065 ns/op
//PASS
//ok      concurrency_in_go/sync  50.802s
