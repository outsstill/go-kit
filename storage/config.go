package storage

type Config struct {
	Prefix string `mapstructure:"prefix" yaml:"prefix"`
	Driver string `mapstructure:"driver" yaml:"driver"` // local | s3 | minio

	SizeLimit int64    `mapstructure:"size_limit" yaml:"size_limit"`
	Ext       []string `mapstructure:"ext" yaml:"ext"`

	Oss   OssConfig   `mapstructure:"oss" yaml:"oss"`
	Local LocalConfig `mapstructure:"local" yaml:"local"`
	S3    S3Config    `mapstructure:"s3" yaml:"s3"`
	Minio MinioConfig `mapstructure:"minio" yaml:"minio"`
}

type LocalConfig struct {
	BasePath     string `mapstructure:"base_path"`
	BaseURL      string `mapstructure:"base_url"`
	StaticPrefix string `mapstructure:"static_prefix"` //访问的替换字段避免暴露真实路径
}

type OssConfig struct {
	Region   string `mapstructure:"region"`
	Bucket   string `mapstructure:"bucket"`
	Key      string `mapstructure:"key"`
	Secret   string `mapstructure:"secret"`
	Domain   string `mapstructure:"domain"`
	RoleArn  string `mapstructure:"role_arn"`
	Duration int64  `mapstructure:"duration"`
	Endpoint string `mapstructure:"endpoint"`
}

type S3Config struct {
	Endpoint  string `mapstructure:"endpoint"`
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Region    string `mapstructure:"region"`
}

type MinioConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	SSL       bool   `mapstructure:"ssl"`
}
