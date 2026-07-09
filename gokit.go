package gokit

import (
	"context"

	"github.com/outsstill/go-kit/cache"
	"github.com/outsstill/go-kit/captcha"
	"github.com/outsstill/go-kit/config"
	"github.com/outsstill/go-kit/database"
	"github.com/outsstill/go-kit/database/mysql"
	"github.com/outsstill/go-kit/logger"
	"github.com/outsstill/go-kit/redis"
	"github.com/outsstill/go-kit/storage"
)

type App struct {
	Config  *config.Config
	DB      database.Database
	Redis   *redis.RedisClient
	Cache   cache.Cache
	Captcha *captcha.Captcha
	Storage storage.IStorage
	Logger  *logger.Logger
}

func Bootstrap(configPath string) (*App, error) {

	app := &App{}

	// config
	configObj, err := config.LoadConfig(configPath)

	if err != nil {
		return nil, err
	}
	app.Config = configObj

	// log
	err = logger.Init(configObj.Logger.ToLoggerConfig())

	if err != nil {
		return nil, err
	}

	app.Logger = logger.LogDefault

	// db
	db, err := mysql.New("default", configObj.DB.ToMySQL())

	if err != nil {
		return nil, err
	}

	app.DB = db

	// redis
	redisClient, err := redis.New(configObj.Redis.ToRedis(), context.Background())
	if err != nil {
		return nil, err
	}

	app.Redis = redisClient

	// cache
	redisCache, err := redis.NewCache(configObj.Redis.ToRedis(), context.Background())
	if err != nil {
		return nil, err
	}
	cacheObj := cache.NewRedisCache(redisCache.Client, nil)

	app.Cache = cacheObj

	// 验证码
	captchaObj, err := captcha.NewCaptcha(app.Redis.Client, configObj.Captcha.ToCaptcha(), nil)
	if err != nil {
		return nil, err
	}

	app.Captcha = captchaObj

	// storage
	storageManager, err := storage.New(configObj.Storage.ToStorage())
	if err != nil {
		return nil, err
	}

	app.Storage = storageManager.Driver()

	return app, nil
}
