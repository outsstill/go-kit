package postgres

import (
	"context"
	"database/sql"

	"github.com/outsstill/go-kit/database"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	name string
	db   *gorm.DB
	sql  *sql.DB
}

func New(
	name string,
	cfg database.Config,
) (*Client, error) {

	gdb, err := gorm.Open(
		postgres.Open(cfg.DSN),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}

	return &Client{
		name: name,
		db:   gdb,
		sql:  sqlDB,
	}, nil
}

func (c *Client) Name() string {
	return c.name
}

func (c *Client) DB() *gorm.DB {
	return c.db
}

func (c *Client) Ping(ctx context.Context) error {
	return c.sql.PingContext(ctx)
}

func (c *Client) Close() error {
	return c.sql.Close()
}
