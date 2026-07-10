package captcha

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Length          int     `mapstructure:"length" json:"length"`
	Width           int     `mapstructure:"width" json:"width"`
	Height          int     `mapstructure:"height" json:"height"`
	DotCount        int     `mapstructure:"dot_count" json:"dot_count"`
	UseNumber       bool    `mapstructure:"use_number" json:"use_number"`
	Expiration      int64   `mapstructure:"expiration" json:"expiration"`
	Prefix          string  `mapstructure:"prefix" json:"prefix"`
	ClearOnVerify   bool    `mapstructure:"clear_on_verify" json:"clear_on_verify"`
	Charset         string  `mapstructure:"charset" json:"charset"`
	Maxskew         float64 `mapstructure:"maxskew" json:"maxskew"`
	ShowLineOptions int     `mapstructure:"show_line_options" json:"show_line_options"`
	TestingKey      string  `mapstructure:"testing_key" json:"testing_key"`
	DebugExpireTime int64   `mapstructure:"debug_expire_time" json:"debug_expire_time"`
}

var bgCtx = context.Background()

type Captcha struct {
	store  *redisStore
	driver base64Captcha.Driver
	cfg    Config
}

type redisStore struct {
	rdb *redis.Client
	cfg Config
}

func NewCaptcha(rdb *redis.Client, cfg Config, ctx context.Context) (*Captcha, error) {

	if rdb == nil {
		return nil, errors.New("redis client is nil")
	}

	if ctx == nil {
		ctx = bgCtx
	}

	// default config
	if cfg.Length <= 0 {
		cfg.Length = 4
	}
	if cfg.Width <= 0 {
		cfg.Width = 240
	}
	if cfg.Height <= 0 {
		cfg.Height = 80
	}
	if cfg.DotCount <= 0 {
		cfg.DotCount = 80
	}
	if cfg.Expiration <= 0 {
		cfg.Expiration = int64(5 * time.Minute)
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "captcha:"
	}

	store := &redisStore{
		rdb: rdb,
		cfg: cfg,
	}

	var driver base64Captcha.Driver

	if cfg.UseNumber {
		driver = base64Captcha.NewDriverDigit(
			cfg.Height,
			cfg.Width,
			cfg.Length,
			cfg.Maxskew,
			cfg.DotCount,
		)
	} else {
		if cfg.Charset == "" {
			cfg.Charset = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		}
		driver = base64Captcha.NewDriverString(
			cfg.Height,
			cfg.Width,
			cfg.DotCount,
			cfg.ShowLineOptions,
			cfg.Length,
			cfg.Charset,
			nil,
			nil,
			nil,
		)
	}

	return &Captcha{
		store:  store,
		driver: driver,
		cfg:    cfg,
	}, nil
}

/*
========================
	Generate
========================
*/

func (c *Captcha) Generate() (id string, b64s, answer string, err error) {

	cap := base64Captcha.NewCaptcha(c.driver, c.store)
	return cap.Generate()
}

func (c *Captcha) Verify(id string, code string) (bool, error) {

	if id == "" || code == "" {
		return false, errors.New("invalid captcha params")
	}

	return c.store.Verify(id, code, c.cfg.ClearOnVerify), nil
}

func (c *Captcha) Reload(id string) (nid, b64s, answer string, err error) {

	if id == "" {
		return "", "", "", errors.New("empty captcha id")
	}

	_ = c.store.Delete(id)

	return c.Generate()
}

func (c *Captcha) TTL(id string) (time.Duration, error) {

	if id == "" {
		return 0, errors.New("empty captcha id")
	}

	return c.store.TTL(id)
}

func (c *Captcha) Delete(id string) error {

	if id == "" {
		return errors.New("empty captcha id")
	}

	return c.store.Delete(id)
}

func (s *redisStore) Set(id string, value string) error {
	err := s.rdb.Set(bgCtx, s.cfg.Prefix+id, value, time.Duration(s.cfg.Expiration)).Err()

	return err
}

func (s *redisStore) Get(id string, clear bool) string {

	key := s.cfg.Prefix + id

	val, err := s.rdb.Get(bgCtx, key).Result()
	if err != nil {
		return ""
	}

	if clear && s.cfg.ClearOnVerify {
		_ = s.rdb.Del(bgCtx, key).Err()
	}

	return val
}

func (s *redisStore) Verify(id, answer string, clear bool) bool {
	return strings.EqualFold(
		s.Get(id, clear),
		answer,
	)
}

func (s *redisStore) Delete(id string) error {
	return s.rdb.Del(bgCtx, s.cfg.Prefix+id).Err()
}

func (s *redisStore) TTL(id string) (time.Duration, error) {
	return s.rdb.TTL(bgCtx, s.cfg.Prefix+id).Result()
}

// base64Captcha 会调用这个接口
func (s *redisStore) GetFromStore(key string, clear bool) string {
	return s.Get(key, clear)
}

func (s *redisStore) SetToStore(key string, value string) {
	s.Set(key, value)
}
