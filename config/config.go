package config

import (
	"errors"

	"github.com/outsstill/go-kit/helpers"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Config struct {
	App     AppConfig     `mapstructure:"app" yaml:"app"`
	DB      DBConfig      `mapstructure:"db" yaml:"db"`
	Redis   RedisConfig   `mapstructure:"redis" yaml:"redis"`
	Logger  LoggerConfig  `mapstructure:"logger" yaml:"logger"`
	JWT     JWTConfig     `mapstructure:"jwt" yaml:"jwt"`
	Captcha CaptchaConfig `mapstructure:"captcha" yaml:"captcha"`
	Storage StorageConfig `mapstructure:"storage" yaml:"storage"`
	Limit   LimitConfig   `mapstructure:"limit" yaml:"limit"`
	Paging  PagingConfig  `mapstructure:"paging" yaml:"paging"`
	v       *viper.Viper  `mapstructure:"-"`
}

func LoadConfig(path string) (*Config, error) {
	if len(path) == 0 {
		return nil, errors.New("config file path is empty")
	}
	v := viper.New()

	s := &Config{}

	if len(path) > 0 {
		v.SetConfigType("yaml") // 类型
		v.SetConfigFile(path)

		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	if err := v.Unmarshal(s); err != nil {
		return nil, err
	}

	s.v = v

	return s, nil
}

func (c *Config) Get(key string, defaultValue ...any) any {
	return c.internalGet(key, defaultValue...)
}

func (c *Config) GetString(key string, defaultValue ...any) string {
	return cast.ToString(c.internalGet(key, defaultValue...))
}

func (c *Config) GetStringSlice(key string, defaultValue ...any) []string {
	return cast.ToStringSlice(c.internalGet(key, defaultValue...))
}

func (c *Config) GetBool(key string, defaultValue ...any) bool {
	return cast.ToBool(c.internalGet(key, defaultValue...))
}

func (c *Config) GetInt(key string, defaultValue ...any) int {
	return cast.ToInt(c.internalGet(key, defaultValue...))
}

func (c *Config) All() map[string]any {
	return c.v.AllSettings()
}

func (c *Config) GetInt64(path string, defaultValue ...any) int64 {
	return cast.ToInt64(c.internalGet(path, defaultValue...))
}

func (c *Config) GetFloat64(path string, defaultValue ...any) float64 {
	return cast.ToFloat64(c.internalGet(path, defaultValue...))
}

func (c *Config) internalGet(path string, defaultValue ...any) any {

	v := c.v.Get(path)
	// config 或者环境变量不存在的情况
	if !c.v.IsSet(path) || helpers.Empty(v) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return v
}
