package config

import "github.com/outsstill/go-kit/logger"

type LoggerConfig struct {
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

func (c LoggerConfig) ToLoggerConfig() logger.Config {
	return logger.Config{
		Level:      c.Level,
		Filename:   c.Filename,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		Compress:   c.Compress,
		Type:       c.Type,
		Encoding:   c.Encoding,
		ToConsole:  c.ToConsole,
		ToFile:     c.ToFile,
	}
}
