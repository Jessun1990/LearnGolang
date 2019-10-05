package chapter1

import "testing"

/*
go test -v -count=1 ./chapter1 -run TestRaceCondCase
*/
func TestRaceCondCase(t *testing.T) {
	raceCondCase()
}

/*
go test -v -count=1 ./chapter1 -run TestMemAccessSyncCase
*/
func TestMemAccessSyncCase(t *testing.T) {
	memAccessSyncCase()
}

/*
go test -v -count=1 ./chapter1 -run TestMemAccessSyncMtxCase
*/
func TestMemAccessSyncMtxCase(t *testing.T) {
	memAccessSyncMtxCase()
}

/*
go test -v -count=1 ./chapter1 -run TestDeadLockCase
*/
func TestDeadLockCase(t *testing.T) {
	deadLockCase()
}

/*
go test -count=1 -v ./chapter1 -run TestLiveLockCase
*/
func TestLiveLockCase(t *testing.T) {
	liveLockCase()
}

/*
go test -count=1 -v ./chapter1 -run TestStarvationCase
*/
func TestStarvationCase(t *testing.T) {
	starvationCase()
}
