package cache

import "time"

// Cache 缓存接口（生产级：显式错误）
type Cache interface {
	Set(key string, value string, expiration time.Duration) error

	Get(key string) (string, error)

	Has(key string) (bool, error)

	Forget(key string) error

	Forever(key string, value string) error

	Flush() error

	IsAlive() error

	// Increment
	// 1个参数：key +1
	// 2个参数：key + value
	Increment(parameters ...interface{}) (int64, error)

	// Decrement
	// 1个参数：key -1
	// 2个参数：key - value
	Decrement(parameters ...interface{}) (int64, error)
}
