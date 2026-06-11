package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/outsstill/go-kit/database"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Client struct {
	name string

	db  *gorm.DB
	sql *sql.DB
}

func New(
	name string,
	cfg database.Config,
) (*Client, error) {

	dsn := cfg.DSN

	if dsn == "" {
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)
	}

	gdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime))

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
