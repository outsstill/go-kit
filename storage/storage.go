package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type IStorage interface {
	// 上传
	Put(ctx context.Context, req *UploadRequest) (*FileObj, error)

	// 下载
	Get(ctx context.Context, key string) (*FileObj, io.ReadCloser, error)

	// 删除
	Delete(ctx context.Context, key string) error

	Driver() string

	URL(ctx context.Context, key string) string
}

type UploadRequest struct {
	SourceType  string
	Path        string
	Filename    string
	ContentType string

	Size   int64
	Reader io.Reader
	// 供直传/分片扩展使用
	Meta map[string]string
}

type FileObj struct {
	ID           string            `json:"id"`
	Bucket       string            `json:"bucket"`
	Path         string            `json:"path"`        // 包含原文件名的上传路径 例: /a/b/c/123.jpg
	Key          string            `json:"key"`         // 包含文件名的完整储存路径 例: /2000/02/01/xxxxxx.jpg
	OriginName   string            `json:"origin_name"` // 原文件名 例: 123.jpg
	StoredName   string            `json:"stored_name"` // 储存文件名 例: xxxxxx.jpg
	Size         int64             `json:"size"`
	ContentType  string            `json:"content_type"`
	Ext          string            `json:"ext"` // 不含 . 例: jpg
	Hash         string            `json:"hash"`
	ETag         string            `json:"e_tag"`
	URL          string            `json:"url"` // 公开或带签名的可访问 URL（可选）
	LastModified time.Time         `json:"last_modified"`
	Driver       string            `json:"driver"`
	IsPublic     bool              `json:"is_public"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// 获取文件存储名称(包含完整路径)
func GetFileStorageRealPath(fileName string, isOriginName bool, prefix string) (string, string) {
	originFileName := fileName
	if !isOriginName {
		fileOriExt := filepath.Ext(fileName) // 获取文件扩展名 这里包含了 .
		//randomNumber := app.GetRandomNumber(16)
		randomNumber := uuid.NewString()
		// fileNameNoExt := fileName[:len(fileName)-len(fileOriExt)] // 文件名称 不含 .和后缀
		originFileName = cast.ToString(randomNumber) + fileOriExt
	}

	objectName := GetFileStoragePathPrefix(prefix) + "/" + originFileName

	return objectName, originFileName
}

func GetFileStoragePathPrefix(prefix string) string {
	return fmt.Sprintf("%s/%s", prefix, time.Now().Format("20060102"))
}

func safeJoin(base, p string) (string, error) {
	base = filepath.Clean(base)

	full := filepath.Join(base, p)
	full = filepath.Clean(full)

	rel, err := filepath.Rel(base, full)
	if err != nil {
		return "", err
	}

	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes base directory")
	}

	return full, nil
}

func GetFileExt(fileName string) string {
	fileOriExt := filepath.Ext(fileName) // 获取文件扩展名 这里包含了 .

	if strings.HasPrefix(fileOriExt, ".") {
		fileOriExt = fileOriExt[1:]
	}

	fileOriExt = strings.ToLower(fileOriExt)

	return fileOriExt
}
