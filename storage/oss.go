package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts "github.com/alibabacloud-go/sts-20150401/v2/client"
	"github.com/alibabacloud-go/tea/tea"
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
	objectName = strings.Trim(objectName, "/")
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
		Key:          objectName,
		StoredName:   nowFileName,
		OriginName:   fileName,
		Path:         objectName,
		Ext:          GetFileExt(fileName),
		Driver:       l.Driver(),
		Size:         fileSize,
		MimeType:     in.MimeType,
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

func (l *OssStorage) Certificate(ctx context.Context, in *UploadRequest) (*UploadCredential, error) {

	config := &openapi.Config{
		AccessKeyId:     tea.String(l.cfg.Oss.Key),
		AccessKeySecret: tea.String(l.cfg.Oss.Secret),
		Endpoint:        tea.String(l.cfg.Oss.EndpointSts),
	}

	client, err := sts.NewClient(config)
	if err != nil {
		return nil, err
	}

	req := &sts.AssumeRoleRequest{
		RoleArn:         tea.String(l.cfg.Oss.RoleArn),
		RoleSessionName: tea.String(in.SessionName),
		DurationSeconds: tea.Int64(l.cfg.Oss.Duration),
	}

	resp, err := client.AssumeRole(req)
	if err != nil {
		return nil, err
	}

	expire, _ := time.Parse(
		time.RFC3339,
		tea.StringValue(resp.Body.Credentials.Expiration),
	)

	files := in.Files

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found")
	}

	returnFiles := make([]UploadFileKey, 0, len(files))

	for _, file := range files {
		if file.Filename != "" && file.UUID != "" {
			key, _ := GetFileStorageRealPath(file.Filename, false, l.cfg.Prefix)
			returnFiles = append(returnFiles, UploadFileKey{
				Filename: file.Filename,
				Key:      key,
				UUID:     file.UUID,
				URL:      l.URL(ctx, key),
			})
		}
	}

	//respJson, _ := json.Marshal(resp)

	cfg := &UploadCredential{
		SourceType:      in.SourceType,
		AccessKeyID:     tea.StringValue(resp.Body.Credentials.AccessKeyId),
		AccessKeySecret: tea.StringValue(resp.Body.Credentials.AccessKeySecret),
		SecurityToken:   tea.StringValue(resp.Body.Credentials.SecurityToken),
		ExpireAt:        expire.Unix(),
		Driver:          l.Driver(),
		Bucket:          l.cfg.Oss.Bucket,
		Region:          l.cfg.Oss.Region,
		Endpoint:        l.cfg.Oss.Endpoint, // 前端上传地址所用
		Files:           returnFiles,
	}

	return cfg, nil
}
