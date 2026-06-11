package config

import "github.com/outsstill/go-kit/captcha"

type CaptchaConfig struct {
	Length          int     `mapstructure:"length" json:"length"`
	Width           int     `mapstructure:"width" json:"width"`
	Height          int     `mapstructure:"height" json:"height"`
	NoiseCount      int     `mapstructure:"noise_count" json:"noise_count"`
	UseNumber       bool    `mapstructure:"use_number" json:"use_number"`
	Expiration      int64   `mapstructure:"expiration" json:"expiration"`
	Prefix          string  `mapstructure:"prefix" json:"prefix"`
	ClearOnVerify   bool    `mapstructure:"clear_on_verify" json:"clear_on_verify"`
	Charset         string  `mapstructure:"charset" json:"charset"`
	Maxskew         float64 `mapstructure:"maxskew" json:"maxskew"`
	ShowLineOptions int     `mapstructure:"show_line_options" json:"show_line_options"`
}

func (c *CaptchaConfig) ToCaptcha() captcha.Config {
	return captcha.Config{
		Length:          c.Length,
		Width:           c.Width,
		Height:          c.Height,
		NoiseCount:      c.NoiseCount,
		UseNumber:       c.UseNumber,
		Expiration:      c.Expiration,
		Prefix:          c.Prefix,
		ClearOnVerify:   c.ClearOnVerify,
		Charset:         c.Charset,
		Maxskew:         c.Maxskew,
		ShowLineOptions: c.ShowLineOptions,
	}
}
