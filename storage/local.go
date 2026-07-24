package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type LocalStorage struct {
	cfg Config
}

func NewLocal(cfg Config) (*LocalStorage, error) {
	return &LocalStorage{
		cfg: cfg,
	}, nil
}

func (l *LocalStorage) Driver() string { return "local" }

func (l *LocalStorage) Put(ctx context.Context, in *UploadRequest) (*FileObj, error) {

	f := &FileObj{}
	realPath, nowFileName := GetFileStorageRealPath(in.Filename, false, l.cfg.Prefix)

	fullPath, err := safeJoin(l.cfg.Local.BasePath, realPath) // 真实的地址

	if err != nil {
		return nil, err
	}

	fullPath = strings.ReplaceAll(fullPath, "\\", "/")

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	size, err := io.Copy(file, in.Reader)
	if err != nil {
		return nil, err
	}

	f.Driver = l.Driver()
	f.Size = size
	f.ID = uuid.NewString()
	f.ContentType = in.ContentType
	f.MimeType = in.MimeType
	f.Path = fullPath
	f.Key = strings.Trim(realPath, "/")
	f.Ext = GetFileExt(in.Filename)
	f.OriginName = in.Filename
	f.StoredName = nowFileName
	f.URL = l.URL(ctx, fullPath)
	f.LastModified = time.Now()

	return f, nil
}

func (l *LocalStorage) Get(ctx context.Context, key string) (*FileObj, io.ReadCloser, error) {
	full, err := safeJoin(l.cfg.Local.BasePath, key)
	if err != nil {
		return &FileObj{}, nil, err
	}
	f, err := os.Open(full)
	if err != nil {
		return &FileObj{}, nil, err
	}
	st, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, nil, err
	}
	return &FileObj{Key: key, Size: st.Size(), LastModified: st.ModTime()}, f, nil
}

func (l *LocalStorage) Delete(ctx context.Context, key string) error {
	full, err := safeJoin(l.cfg.Local.BasePath, key)
	if err != nil {
		return err
	}

	err = os.Remove(full)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func (l *LocalStorage) URL(ctx context.Context, path string) string {
	path = strings.ReplaceAll(path, l.cfg.Local.BasePath, l.cfg.Local.StaticPrefix)
	return strings.TrimRight(l.cfg.Local.BaseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

func (l *LocalStorage) Certificate(ctx context.Context, in *UploadRequest) (*UploadCredential, error) {
	return &UploadCredential{}, nil
}
