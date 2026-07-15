package storage

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/spf13/cast"
)

type OssStorage struct {
	cfg       Config
	ossClient *oss.Client
}

func NewOssStorage(cfg Config) (*OssStorage, error) {
	var (
		//region     = "oss-cn-shanghai.aliyuncs.com"
		region = cast.ToString(cfg.Oss.Region)
	)
	_ = os.Setenv("OSS_ACCESS_KEY_ID", cast.ToString(cfg.Oss.Key))
	_ = os.Setenv("OSS_ACCESS_KEY_SECRET", cast.ToString(cfg.Oss.Secret))
	ossCfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)
	client := oss.NewClient(ossCfg)
	return &OssStorage{
		cfg:       cfg,
		ossClient: client,
	}, nil
}

func (l *OssStorage) Driver() string { return "oss" }

func (l *OssStorage) Put(ctx context.Context, in *UploadRequest) (*FileObj, error) {

	// 从请求中获取文件
	fileName := in.Filename
	fileSize := in.Size
	objectName, nowFileName := GetFileStorageRealPath(fileName, false, l.cfg.Prefix)

	res, err := l.ossClient.PutObject(ctx, &oss.PutObjectRequest{
		Bucket:      oss.Ptr(l.cfg.Oss.Bucket),
		Key:         oss.Ptr(objectName),
		Body:        in.Reader,
		ContentType: oss.Ptr(in.ContentType),
	})

	if err != nil {
		return nil, err
	}

	info := &FileObj{
		Bucket:       l.cfg.Oss.Bucket,
		Key:          strings.Trim(objectName, "/"),
		StoredName:   nowFileName,
		OriginName:   fileName,
		Path:         in.Path,
		Ext:          GetFileExt(fileName),
		Driver:       l.Driver(),
		Size:         fileSize,
		ContentType:  in.ContentType,
		ETag:         oss.ToString(res.ETag), // 可计算 md5
		LastModified: time.Now(),
		URL:          l.URL(ctx, objectName),
	}

	return info, nil
}

func (l *OssStorage) Get(ctx context.Context, key string) (*FileObj, io.ReadCloser, error) {

	resp, err := l.ossClient.GetObject(
		ctx,
		&oss.GetObjectRequest{
			Bucket: oss.Ptr(l.cfg.Oss.Bucket),
			Key:    oss.Ptr(key),
		},
	)

	if err != nil {
		return nil, nil, err
	}

	info := &FileObj{
		Bucket: l.cfg.Oss.Bucket,
		Key:    key,
		Driver: l.Driver(),
		URL:    strings.TrimRight(l.cfg.Oss.Domain, "/") + "/" + key,
	}

	return info, resp.Body, nil
}

func (l *OssStorage) Delete(ctx context.Context, key string) error {

	_, err := l.ossClient.DeleteObject(ctx, &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(l.cfg.Oss.Bucket),
		Key:    oss.Ptr(key),
	})

	if err != nil {
		return err
	}

	return nil
}

func (l *OssStorage) URL(ctx context.Context, key string) string {
	return strings.TrimRight(l.cfg.Oss.Domain, "/") + "/" + key
}
