// Package logger 处理日志相关逻辑
package logger

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level      string `mapstructure:"level" yaml:"level"`
	Filename   string `mapstructure:"filename" yaml:"filename"`
	MaxSize    int    `mapstructure:"max_size" yaml:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" yaml:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" yaml:"max_age"`
	Compress   bool   `mapstructure:"compress" yaml:"compress"`
	Type       string `mapstructure:"type" yaml:"type"`         // daily / normal
	Encoding   string `mapstructure:"encoding" yaml:"encoding"` // json / console
	ToConsole  bool   `mapstructure:"to_console" yaml:"to_console"`
	ToFile     bool   `mapstructure:"to_file" yaml:"to_file"`
}

type Logger struct {
	z *zap.Logger
}

// New 初始化 logger
func New(cfg Config) *Logger {

	writeSyncer := getWriteSyncer(cfg)

	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	core := zapcore.NewCore(
		getEncoder(cfg),
		writeSyncer,
		level,
	)

	z := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return &Logger{z: z}
}

func (l *Logger) Zap() *zap.Logger {
	return l.z
}

func (l *Logger) Sync() error {
	return l.z.Sync()
}

// =========================
// core API（统一风格）
// =========================

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.z.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.z.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.z.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.z.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.z.Fatal(msg, fields...)
}

// =========================
// 工具方法
// =========================

// With 增加上下文字段（非常重要）
func (l *Logger) Named(name string) *Logger {
	return &Logger{
		z: l.z.Named(name),
	}
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		z: l.z.With(fields...),
	}
}

// Dump 调试对象
func (l *Logger) Dump(v any) {
	l.z.Warn("dump", zap.Any("data", v))
}

// =========================
// encoder
// =========================

func getEncoder(cfg Config) zapcore.Encoder {
	ec := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,      // 每行日志的结尾添加 "\n"
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 日志级别名称大写，如 ERROR、INFO
		EncodeTime:     customTimeEncoder,              // 时间格式，我们自定义为 2006-01-02 15:04:05
		EncodeDuration: zapcore.SecondsDurationEncoder, // 执行时间，以秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,     // Caller 短格式，如：types/converter.go:17，长格式为绝对路径
	}

	switch cfg.Encoding {
	case "console":
		ec.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(ec)
	default:
		return zapcore.NewJSONEncoder(ec)
	}
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// =========================
// writer
// =========================

func getWriteSyncer(cfg Config) zapcore.WriteSyncer {

	filename := cfg.Filename

	// daily log
	if cfg.Type == "daily" {
		logname := time.Now().Format("2006-01-02.log")
		filename = strings.ReplaceAll(filename, "logs.log", logname)
	}

	lumber := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	switch {
	case cfg.ToConsole && cfg.ToFile:
		return zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(lumber),
		)
	case cfg.ToConsole:
		return zapcore.AddSync(os.Stdout)
	default:
		return zapcore.AddSync(lumber)
	}
}
