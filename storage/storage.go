package storage

import "context"

type Storage interface {
	// 上传
	Put(ctx context.Context, path string, data []byte, opts ...Option) (string, error)

	// 下载
	Get(ctx context.Context, path string) ([]byte, error)

	// 删除
	Delete(ctx context.Context, path string) error

	// 获取访问URL（可选）
	URL(ctx context.Context, path string) (string, error)
}
