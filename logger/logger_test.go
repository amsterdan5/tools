package logger

import "testing"

func TestLog(t *testing.T) {
	InitLogger(LogModFile, "")

	Error("err")
	Info("info")

	SetTraceId(0, genTraceId())

	Info("again")
}
