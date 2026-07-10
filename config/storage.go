package config

import "github.com/outsstill/go-kit/storage"

type StorageConfig struct {
	Driver    string   `mapstructure:"driver" yaml:"driver"` // local | s3 | minio
	SizeLimit int64    `mapstructure:"size_limit" yaml:"size_limit"`
	Ext       []string `mapstructure:"ext" yaml:"ext"`

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
	Region     string `mapstructure:"region"`
	BucketName string `mapstructure:"bucket_name"`
	Key        string `mapstructure:"key"`
	Secret     string `mapstructure:"secret"`
	Domain     string `mapstructure:"domain"`
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

func (c StorageConfig) ToStorage() storage.Config {
	return storage.Config{
		Driver:    c.Driver,
		SizeLimit: c.SizeLimit,
		Ext:       c.Ext,
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
