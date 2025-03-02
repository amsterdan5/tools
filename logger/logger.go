package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggerClient *zap.Logger

func InitLogger(mod LogMod, logPath string) {
	var ws zapcore.WriteSyncer
	switch mod {
	case LogModFile: // 文件
		ws = getWriter(logPath)
	case LogModConsole: // 终端
		ws = zapcore.AddSync(os.Stdout)
	case LogModBoth: // 文件+终端
		ws = zapcore.NewMultiWriteSyncer(getWriter(logPath), zapcore.AddSync(os.Stdout))
	default:
		ws = getWriter(logPath)
	}

	// 实例化
	loggerClient = zap.New(zapcore.NewCore(getEncoder(), ws, zapcore.DebugLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		// zap.AddStacktrace(zap.ErrorLevel),
	)
}

func getEncoder() zapcore.Encoder {
	encoder := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "myLogger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewConsoleEncoder(encoder)
}

func getWriter(logPath string) zapcore.WriteSyncer {
	if logPath == "" {
		logPath = defaultLogPath
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,    // 在进行切割之前，日志文件的最大大小(MB)
		MaxBackups: 5,     // 保留旧文件的最大个数
		MaxAge:     30,    // 保留旧文件的最大天数
		Compress:   false, // 是否压缩 / 归档旧文件
	}

	return zapcore.AddSync(lumberJackLogger)
}

// debug
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// 普通消息
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// 提醒
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// 错误
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// debug
func Debugf(msg string, fields ...interface{}) {
	GetLogger().Sugar().Debugf(msg, fields...)
}

// 普通消息
func Infof(msg string, fields ...interface{}) {
	GetLogger().Sugar().Infof(msg, fields...)
}

// 提醒
func Warnf(msg string, fields ...interface{}) {
	GetLogger().Sugar().Warnf(msg, fields...)
}

// 错误
func Errorf(msg string, fields ...interface{}) {
	GetLogger().Sugar().Errorf(msg, fields...)
}

func Sync() error {
	return loggerClient.Sync()
}

func WithOptions(opts ...zap.Option) *zap.Logger {
	return loggerClient.WithOptions(opts...)
}

func GetLogger() *zap.Logger {
	gid, traceId := GetTraceId()
	return loggerClient.With(zap.String("traceId", traceId), zap.Int64("gid", gid))
}
