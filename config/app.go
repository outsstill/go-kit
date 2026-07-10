package config

type AppConfig struct {
	Name     string `mapstructure:"name" yaml:"name"`
	Key      string `mapstructure:"key" yaml:"key"`
	Url      string `mapstructure:"url" yaml:"url"`
	HttpPort string `mapstructure:"http_port" yaml:"http_port"`
	FileUrl  string `mapstructure:"file_url" yaml:"file_url"`
	Env      string `mapstructure:"env" yaml:"env"`
	Version  string `mapstructure:"version" yaml:"version"`
	Debug    bool   `mapstructure:"debug" yaml:"debug"`
	Timezone string `mapstructure:"timezone" yaml:"timezone"`
}
