package storage

import "errors"

type Manager struct {
	driver Storage
}

func New(cfg Config) (*Manager, error) {
	var s Storage
	var err error

	switch cfg.Driver {
	case "local":
		s, err = NewLocal(cfg.Local)
	case "s3":
		s, err = NewS3(cfg.S3)
	case "minio":
		s, err = NewMinio(cfg.Minio)
	default:
		return nil, errors.New("unknown storage driver")
	}

	if err != nil {
		return nil, err
	}

	return &Manager{driver: s}, nil
}

func (m *Manager) Driver() Storage {
	return m.driver
}
