package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client  *redis.Client
	context context.Context
}

func NewRedisCache(client *redis.Client, ctx context.Context) *RedisCache {
	if ctx == nil {
		ctx = context.Background()
	}

	return &RedisCache{
		client:  client,
		context: ctx,
	}
}

func (r *RedisCache) Set(key string, value string, expiration time.Duration) error {
	return r.client.Set(r.context, key, value, expiration).Err()
}

func (r *RedisCache) Get(key string) (string, error) {
	val, err := r.client.Get(r.context, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}

	return val, nil
}

func (r *RedisCache) Has(key string) (bool, error) {
	n, err := r.client.Exists(r.context, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *RedisCache) Forget(key string) error {
	return r.client.Del(r.context, key).Err()
}

func (r *RedisCache) Forever(key string, value string) error {
	return r.client.Set(r.context, key, value, 0).Err()
}

func (r *RedisCache) Flush() error {
	return r.client.FlushDB(r.context).Err()
}

func (r *RedisCache) IsAlive() error {
	return r.client.Ping(r.context).Err()
}

func (r *RedisCache) Increment(parameters ...interface{}) (int64, error) {

	if len(parameters) == 0 {
		return 0, errors.New("missing key")
	}

	key, ok := parameters[0].(string)
	if !ok {
		return 0, errors.New("key must be string")
	}

	if len(parameters) == 1 {
		return r.client.Incr(r.context, key).Result()
	}

	val, ok := parameters[1].(int64)
	if !ok {
		return 0, errors.New("increment value must be int64")
	}

	return r.client.IncrBy(r.context, key, val).Result()
}

func (r *RedisCache) Decrement(parameters ...interface{}) (int64, error) {

	if len(parameters) == 0 {
		return 0, errors.New("missing key")
	}

	key, ok := parameters[0].(string)
	if !ok {
		return 0, errors.New("key must be string")
	}

	if len(parameters) == 1 {
		return r.client.Decr(r.context, key).Result()
	}

	val, ok := parameters[1].(int64)
	if !ok {
		return 0, errors.New("decrement value must be int64")
	}

	return r.client.DecrBy(r.context, key, val).Result()
}
