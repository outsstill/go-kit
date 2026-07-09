package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/outsstill/go-kit/logger"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	Username     string        `mapstructure:"username"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	CacheDB      int           `mapstructure:"cache_db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	MaxRetries   int           `mapstructure:"max_retries"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type RedisClient struct {
	Client  *redis.Client
	Context context.Context
}

func New(cfg Config, ctx context.Context) (*RedisClient, error) {
	return newRedis(cfg, ctx, cfg.DB)
}

func NewCache(cfg Config, ctx context.Context) (*RedisClient, error) {
	return newRedis(cfg, ctx, cfg.CacheDB)
}

func newRedis(cfg Config, ctx context.Context, db int) (*RedisClient, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	rds := &RedisClient{
		Context: ctx,
		Client: redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
			Username:     cfg.Username,
			Password:     cfg.Password,
			DB:           db,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConns,
			MaxRetries:   cfg.MaxRetries,
			DialTimeout:  cfg.DialTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		}),
	}

	if err := rds.Client.Ping(rds.Context).Err(); err != nil {
		logger.ErrorString("redis", "newRedis", err.Error())
		return nil, err
	}

	return rds, nil
}

func (r *RedisClient) Raw() *redis.Client {
	return r.Client
}

func (r *RedisClient) Ping() error {
	return r.Client.Ping(r.Context).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.Client.Get(r.Context, key).Result()
}

func (r *RedisClient) Set(key string, value any, expiration time.Duration) error {
	return r.Client.Set(r.Context, key, value, expiration).Err()
}

func (r *RedisClient) Del(keys ...string) error {
	return r.Client.Del(r.Context, keys...).Err()
}

func (r *RedisClient) Exists(keys ...string) (int64, error) {
	return r.Client.Exists(r.Context, keys...).Result()
}

func (r *RedisClient) Expire(key string, expiration time.Duration) (bool, error) {
	return r.Client.Expire(r.Context, key, expiration).Result()
}

func (r *RedisClient) TTL(key string) (time.Duration, error) {
	return r.Client.TTL(r.Context, key).Result()
}

func (r *RedisClient) Incr(key string) (int64, error) {
	return r.Client.Incr(r.Context, key).Result()
}

func (r *RedisClient) Decr(key string) (int64, error) {
	return r.Client.Decr(r.Context, key).Result()
}

func (r *RedisClient) IncrBy(key string, value int64) (int64, error) {
	return r.Client.IncrBy(r.Context, key, value).Result()
}

func (r *RedisClient) DecrBy(key string, value int64) (int64, error) {
	return r.Client.DecrBy(r.Context, key, value).Result()
}

func (r *RedisClient) HGet(key, field string) (string, error) {
	return r.Client.HGet(r.Context, key, field).Result()
}

func (r *RedisClient) HSet(key string, values ...any) (int64, error) {
	return r.Client.HSet(r.Context, key, values...).Result()
}

func (r *RedisClient) HGetAll(key string) (map[string]string, error) {
	return r.Client.HGetAll(r.Context, key).Result()
}

func (r *RedisClient) LPush(key string, values ...any) (int64, error) {
	return r.Client.LPush(r.Context, key, values...).Result()
}

func (r *RedisClient) RPop(key string) (string, error) {
	return r.Client.RPop(r.Context, key).Result()
}

func (r *RedisClient) LRange(key string, start, stop int64) ([]string, error) {
	return r.Client.LRange(r.Context, key, start, stop).Result()
}

func (r *RedisClient) FlushDB() error {
	return r.Client.FlushDB(r.Context).Err()
}

func (r *RedisClient) FlushAll() error {
	return r.Client.FlushAll(r.Context).Err()
}
