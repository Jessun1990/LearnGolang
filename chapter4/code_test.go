package chapter4

import "testing"

/*
go test ./chapter4 -v -count=1 -run TestCodeExample
*/
func TestCodeExample(t *testing.T) {
	codeExample()
}

/*
go test ./chapter4 -v -count=1 -run TestCodeExample2
*/
func TestCodeExample2(t *testing.T) {
	codeExample2()
}

/*
go test ./chapter4 -v -count=1 -run TestCodeExample3
*/
func TestCodeExample3(t *testing.T) {
	codeExample3()
}

/*
go test ./chapter4 -v -count=1 -run TestGoroutineExample3
*/
func TestForSelectExample6(t *testing.T) {
	forSelectExample6()
}

/*
go test ./chapter4 -v -count=1 -run TestGoroutineExample4
*/
func TestForSelectExample7(t *testing.T) {
	forSelectExample7()
}

/*
go test ./chapter4 -v -count=1 -run TestOrChannelExample
*/
func TestOrChannelExample(t *testing.T) {
	orChannelExample()
}

/*
go test ./chapter4 -v -count=1 -run TestErrHandleExample
*/
func TestErrHandleExample(t *testing.T) {
	errHandleExample()
}

/*
go test ./chapter4 -v -count=1 -run TestErrHandleExample2
*/
func TestErrHandleExample2(t *testing.T) {
	errHandleExample2()
}

/*
go test ./chapter4 -v -count=1 -run TestErrHandleExample3
*/
func TestErrHandleExample3(t *testing.T) {
	errHandleExample3()
}

/*
go test ./chapter4 -v -count=1 -run TestPipelineExample
*/
func TestPipelineExample(t *testing.T) {
	pipelineExample()
}

/*
go test ./chapter4 -v -count=1 -run TestPipelineExample3
*/
func TestPipelineExample3(t *testing.T) {
	pipelineExample3()
}

/*
go test ./chapter4 -v -count=1 -run TestGeneratorExample
*/
func TestGeneratorExample(t *testing.T) {
	generatorExample()
}

/*
go test ./chapter4 -v -count=1 -run TestGeneratorExample2
*/
func TestGeneratorExample2(t *testing.T) {
	generatorExample2()
}

/*
go test ./chapter4 -v -count=1 -run TestForSelectExample8
*/
func TestForSelectExample8(t *testing.T) {
	forSelectExample8()
}

/*
go test ./chapter4 -v -count=1 -run TestFanInFanOutExample
*/
func TestFanInFanOutExample(t *testing.T) {
	fanInFanOutExmaple()
}

/*
go test ./chapter4 -v -count=1 -run TestBridgeChannelExample
*/
func TestBridgeChannelExample(t *testing.T) {
	bridgeChannelExample()
}
