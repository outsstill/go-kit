package config

import (
	"github.com/outsstill/go-kit/database"
)

type DBConfig struct {
	Default string      `mapstructure:"default"`
	Mysql   MySQLConfig `mapstructure:"mysql"`
}

type MySQLConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	Database    string `mapstructure:"database"`
	Charset     string `mapstructure:"charset"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"`
	MaxOpenConn int    `mapstructure:"max_open_conn"`
	MaxLifeTime int64  `mapstructure:"max_life_time"`
}

func (c DBConfig) ToMySQL() database.Config {
	return database.Config{
		Driver:      c.Default,
		Host:        c.Mysql.Host,
		Port:        c.Mysql.Port,
		User:        c.Mysql.Username,
		Password:    c.Mysql.Password,
		Database:    c.Mysql.Database,
		MaxIdleConn: c.Mysql.MaxIdleConn,
		MaxOpenConn: c.Mysql.MaxOpenConn,
		MaxLifeTime: c.Mysql.MaxLifeTime,
	}
}
