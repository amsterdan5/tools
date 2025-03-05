package logger

type LogMod uint32

var defaultLogPath = "/home/yoo/log.log"

const (
	LogModFile    = iota + 1 // 文件
	LogModConsole            // 终端
	LogModBoth               // 文件+终端
)
