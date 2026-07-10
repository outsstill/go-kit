package config

type LimitConfig struct {
	Rate        string `mapstructure:"rate" yaml:"rate"` // 过期时间，单位是分钟，一般不超过两个小时
	TestRate    string `mapstructure:"test_rate" yaml:"test_rate"`
	LoginRate   string `mapstructure:"login_rate" yaml:"login_rate"`
	CaptchaRate string `mapstructure:"captcha_rate" yaml:"captcha_rate"`
	StoreRate   string `mapstructure:"store_rate" yaml:"store_rate"`
	UpdateRate  string `mapstructure:"update_rate" yaml:"update_rate"`
	DeleteRate  string `mapstructure:"delete_rate" yaml:"delete_rate"`
	QueryRate   string `mapstructure:"query_rate" yaml:"query_rate"`
}
