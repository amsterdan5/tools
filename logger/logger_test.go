package logger

import (
	"testing"

	"github.com/petermattis/goid"
)

func TestLog(t *testing.T) {
	InitLogger(LogModConsole, "")

	Error("err")
	Info("info")

	SetTraceId(0, genTraceId())

	Info("again")

	SetTraceId(goid.Get(), "")
	Info("set")
}
