package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	DB      DBConfig      `mapstructure:"db"`
	Redis   RedisConfig   `mapstructure:"redis"`
	Logger  LoggerConfig  `mapstructure:"logger"`
	JWT     JWTConfig     `mapstructure:"jwt"`
	Captcha CaptchaConfig `mapstructure:"captcha"`
	Storage StorageConfig `mapstructure:"storage"`
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

	return s, nil
}
