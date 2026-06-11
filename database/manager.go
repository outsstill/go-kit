package database

import (
	"fmt"
	"sync"
)

type Manager struct {
	mu  sync.RWMutex
	dbs map[string]Database
}

func NewManager() *Manager {
	return &Manager{
		dbs: make(map[string]Database),
	}
}

func (m *Manager) Register(name string, db Database) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.dbs[name] = db
}

func (m *Manager) Get(name string) (Database, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	db, ok := m.dbs[name]
	if !ok {
		return nil, fmt.Errorf("database [%s] not found", name)
	}

	return db, nil
}

func (m *Manager) MustGet(name string) Database {
	db, err := m.Get(name)
	if err != nil {
		panic(err)
	}

	return db
}

func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, db := range m.dbs {
		_ = db.Close()
	}

	return nil
}
