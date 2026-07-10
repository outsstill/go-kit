package config

import "github.com/outsstill/go-kit/jwt"

type JWTConfig struct {
	Name       string       `mapstructure:"name" json:"name"`
	Key        string       `mapstructure:"key" json:"key"`
	MaxRefresh int64        `mapstructure:"max_refresh" json:"max_refresh"`
	Timezone   string       `mapstructure:"timezone" json:"timezone"`
	Expires    int64        `mapstructure:"expires" json:"expires"`
	Type       jwt.JWT_TYPE `mapstructure:"type" json:"type"`
}

func (c JWTConfig) ToJWT() jwt.Config {
	return jwt.Config{
		Name:       c.Name,
		Key:        c.Key,
		MaxRefresh: c.MaxRefresh,
		Timezone:   c.Timezone,
		Expires:    c.Expires,
		Type:       c.Type,
	}
}
