package config

import "github.com/outsstill/go-kit/storage"

type StorageConfig struct {
	Driver string `mapstructure:"driver" yaml:"driver"` // local | s3 | minio

	Local LocalConfig `mapstructure:"local" yaml:"local"`
	Oss   OssConfig   `mapstructure:"oss" yaml:"oss"`
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

func (c StorageConfig) ToStorage() storage.Config {
	return storage.Config{
		Driver: c.Driver,
		Local: storage.LocalConfig{
			BasePath:     c.Local.BasePath,
			BaseURL:      c.Local.BaseURL,
			StaticPrefix: c.Local.StaticPrefix,
		},
		Oss: storage.OssConfig{
			Region:     c.Oss.Region,
			BucketName: c.Oss.BucketName,
			Key:        c.Oss.Key,
			Secret:     c.Oss.Secret,
			Domain:     c.Oss.Domain,
		},
		S3: storage.S3Config{
			Endpoint:  c.S3.Endpoint,
			Bucket:    c.S3.Bucket,
			AccessKey: c.S3.AccessKey,
			SecretKey: c.S3.SecretKey,
			Region:    c.S3.Region,
		},
		Minio: storage.MinioConfig{
			Endpoint:  c.Minio.Endpoint,
			Bucket:    c.Minio.Bucket,
			AccessKey: c.Minio.AccessKey,
			SecretKey: c.Minio.SecretKey,
			SSL:       c.Minio.SSL,
		},
	}
}
