package storage

type Config struct {
	Driver string `mapstructure:"driver" yaml:"driver"` // local | s3 | minio

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
	Region     string
	BucketName string
	Key        string
	Secret     string
	Domain     string
}

type S3Config struct {
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
	Region    string
}

type MinioConfig struct {
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
	SSL       bool
}
