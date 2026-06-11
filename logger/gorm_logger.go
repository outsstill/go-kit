package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger GORM Zap Logger
type GormLogger struct {
	ZapLogger     *zap.Logger
	SlowThreshold time.Duration
	LogLevel      gormlogger.LogLevel
}

// NewGormLogger 初始化
func NewGormLogger(z *zap.Logger) *GormLogger {
	return &GormLogger{
		ZapLogger:     z,
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      gormlogger.Info,
	}
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info
func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}
	l.ZapLogger.Info(fmt.Sprintf(msg, args...))
}

// Warn
func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	l.ZapLogger.Warn(fmt.Sprintf(msg, args...))
}

// Error
func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	l.ZapLogger.Error(fmt.Sprintf(msg, args...))
}

// Trace SQL 主逻辑
func (l *GormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (string, int64),
	err error,
) {
	if l.LogLevel == gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.String("cost", microsecondsStr(elapsed)),
	}

	// 1. error 处理
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if l.LogLevel >= gormlogger.Warn {
				l.ZapLogger.Warn("DB NotFound", fields...)
			}
			return
		}

		fields = append(fields, zap.Error(err))

		if l.LogLevel >= gormlogger.Error {
			l.ZapLogger.Error("DB Error", fields...)
		}
		return
	}

	// 2. slow query
	if l.SlowThreshold > 0 && elapsed > l.SlowThreshold {
		if l.LogLevel >= gormlogger.Warn {
			l.ZapLogger.Warn("DB Slow Query", fields...)
		}
		return
	}

	// 3. normal query
	if l.LogLevel >= gormlogger.Info {
		l.ZapLogger.Debug("DB Query", fields...)
	}
}

// microseconds
func microsecondsStr(d time.Duration) string {
	return fmt.Sprintf("%.3fms", float64(d.Nanoseconds())/1e6)
}
