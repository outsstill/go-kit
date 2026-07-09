// Package logger 处理日志相关逻辑
package logger

import (
	"encoding/json"
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

var LogDefault *Logger

func Init(cfg Config) error {

	log, err := New(cfg)

	LogDefault = &Logger{z: log}

	return err
}

// New 初始化 logger
func New(cfg Config) (*zap.Logger, error) {

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

	return z, nil
}

func Zap() *zap.Logger {
	return LogDefault.z
}

func Sync() error {
	return LogDefault.z.Sync()
}

// =========================
// core API（统一风格）
// =========================
func LogIf(err error) {
	if err != nil {
		LogDefault.z.Error("Error Occurred:", zap.Error(err))
	}
}

func Debug(msg string, fields ...zap.Field) {
	LogDefault.z.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	LogDefault.z.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	LogDefault.z.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	LogDefault.z.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	LogDefault.z.Fatal(msg, fields...)
}

// =========================
// 工具方法
// =========================

// With 增加上下文字段（非常重要）
func Named(name string) *Logger {
	return &Logger{
		z: LogDefault.z.Named(name),
	}
}

func With(fields ...zap.Field) *Logger {
	return &Logger{
		z: LogDefault.z.With(fields...),
	}
}

// Dump 调试对象
func Dump(v any) {
	LogDefault.z.Warn("dump", zap.Any("data", v))
}

// LogWarnIf 当 err != nil 时记录 warning 等级的日志
func LogWarnIf(err error) {
	if err != nil {
		LogDefault.z.Warn("Error Occurred:", zap.Error(err))
	}
}

// LogInfoIf 当 err != nil 时记录 info 等级的日志
func LogInfoIf(err error) {
	if err != nil {
		LogDefault.z.Info("Error Occurred:", zap.Error(err))
	}
}

func DebugString(moduleName, name, msg string) {
	LogDefault.z.Debug(moduleName, zap.String(name, msg))
}

func InfoString(moduleName, name, msg string) {
	LogDefault.z.Info(moduleName, zap.String(name, msg))
}

func WarnString(moduleName, name, msg string) {
	LogDefault.z.Warn(moduleName, zap.String(name, msg))
}

func ErrorString(moduleName, name, msg string) {
	LogDefault.z.Error(moduleName, zap.String(name, msg))
}

func FatalString(moduleName, name, msg string) {
	LogDefault.z.Fatal(moduleName, zap.String(name, msg))
}

// DebugJSON 记录对象类型的 debug 日志，使用 json.Marshal 进行编码。调用示例：
//
//	logger.DebugJSON("Auth", "读取登录用户", auth.CurrentUser())
func DebugJSON(moduleName, name string, value interface{}) {
	LogDefault.z.Debug(moduleName, zap.String(name, jsonString(value)))
}

func InfoJSON(moduleName, name string, value interface{}) {
	LogDefault.z.Info(moduleName, zap.String(name, jsonString(value)))
}

func WarnJSON(moduleName, name string, value interface{}) {
	LogDefault.z.Warn(moduleName, zap.String(name, jsonString(value)))
}

func ErrorJSON(moduleName, name string, value interface{}) {
	LogDefault.z.Error(moduleName, zap.String(name, jsonString(value)))
}

func FatalJSON(moduleName, name string, value interface{}) {
	LogDefault.z.Fatal(moduleName, zap.String(name, jsonString(value)))
}

func jsonString(value interface{}) string {
	b, err := json.Marshal(value)
	if err != nil {
		LogDefault.z.Error("Logger", zap.String("JSON marshal error", err.Error()))
	}
	return string(b)
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
		logname := time.Now().Format("20060102.log")
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
