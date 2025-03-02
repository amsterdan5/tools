package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	gLogger "gorm.io/gorm/logger"
)

var (
	traceStr     = gLogger.Green + "%s " + gLogger.Reset + gLogger.Yellow + "[%.3fms] " + gLogger.Reset + "[rows:%v]"
	traceWarnStr = gLogger.Green + "%s " + gLogger.Reset + gLogger.RedBold + "[%.3fms] " + gLogger.Reset + "[rows:%v]"
)

type sqlLogger struct {
	gLogger.Config
}

func NewSqlLogger(config gLogger.Config) *sqlLogger {
	return &sqlLogger{
		config,
	}
}

func (s *sqlLogger) LogMode(level gLogger.LogLevel) gLogger.Interface {
	s.LogLevel = level
	return s
}

func (s *sqlLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if s.LogLevel >= gLogger.Info {
		loggerClient.Info(msg, zap.Any("data", data))
	}
}

func (s *sqlLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if s.LogLevel >= gLogger.Warn {
		loggerClient.Warn(msg, zap.Any("data", data))
	}
}

func (s *sqlLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if s.LogLevel >= gLogger.Error {
		loggerClient.Error(msg, zap.Any("data", data))
	}
}

func (s *sqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if s.LogLevel <= gLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		if elapsed > s.SlowThreshold {
			Errorf(traceWarnStr, sql, float64(elapsed.Nanoseconds())/1e6, rows)
		} else {
			Errorf(traceStr, sql, float64(elapsed.Nanoseconds())/1e6, rows)
		}
	} else {
		if elapsed > s.SlowThreshold {
			Infof(traceWarnStr, sql, float64(elapsed.Nanoseconds())/1e6, rows)
		} else {
			Infof(traceStr, sql, float64(elapsed.Nanoseconds())/1e6, rows)
		}
	}
}
