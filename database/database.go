package database

import (
	"context"

	"gorm.io/gorm"
)

type Database interface {
	Name() string

	DB() *gorm.DB

	Ping(ctx context.Context) error

	Close() error
}
