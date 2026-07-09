package config

import (
	"time"

	"github.com/outsstill/go-kit/redis"
)

type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	Username     string        `mapstructure:"username"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	CacheDB      int           `mapstructure:"cache_db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	MaxRetries   int           `mapstructure:"max_retries"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

func (c *RedisConfig) ToRedis() redis.Config {
	return redis.Config{
		Host:         c.Host,
		Port:         c.Port,
		Username:     c.Username,
		Password:     c.Password,
		DB:           c.DB,
		CacheDB:      c.CacheDB,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
		MaxRetries:   c.MaxRetries,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
	}
}
