package storage

import "errors"

type Manager struct {
	driver IStorage
}

func New(cfg Config) (*Manager, error) {
	var s IStorage
	var err error

	switch cfg.Driver {
	case "local":
		s, err = NewLocal(cfg.Local)
	case "oss":
		s, err = NewOssStorage(cfg.Oss)
	default:
		return nil, errors.New("unknown storage driver")
	}

	if err != nil {
		return nil, err
	}

	return &Manager{driver: s}, nil
}

func (m *Manager) Driver() IStorage {
	return m.driver
}
