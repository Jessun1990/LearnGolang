package sync

import "testing"

// go test ./sync  -run TestShowWaitGroup -v
func TestShowWaitGroup(t *testing.T) {
	ShowWaitGroup()
}

// go test ./sync  -run TestShowMutex -v
func TestShowMutex(t *testing.T) {
	ShowMutex()
}
