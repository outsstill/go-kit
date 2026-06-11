package database

import "time"

type Config struct {
	Driver string `mapstructure:"driver"`

	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`

	MaxIdleConn int           `mapstructure:"max_idle_conn"`
	MaxOpenConn int           `mapstructure:"max_open_conn"`
	MaxLifeTime time.Duration `mapstructure:"max_life_time"`

	DSN string `mapstructure:"dsn"`
}
