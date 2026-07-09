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
	if err := app.LoadConfig(configPath); err != nil {
		return nil, err
	}

	// log
	if err := app.InitLogger(); err != nil {
		return nil, err
	}

	// db
	if err := app.InitDB(); err != nil {
		return nil, err
	}

	// redis
	if err := app.InitRedis(); err != nil {
		return nil, err
	}

	// cache
	if err := app.InitCache(); err != nil {
		return nil, err
	}

	// 验证码
	if err := app.InitCaptcha(); err != nil {
		return nil, err
	}

	// storage
	if err := app.InitStorage(); err != nil {
		return nil, err
	}

	return app, nil
}

func New(configPath string) (*App, error) {
	app := &App{}

	// config
	if err := app.LoadConfig(configPath); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) LoadConfig(configPath string) error {
	configObj, err := config.LoadConfig(configPath)

	if err != nil {
		return err
	}
	a.Config = configObj
	return nil
}

func (a *App) InitLogger() error {
	err := logger.Init(a.Config.Logger.ToLoggerConfig())
	if err != nil {
		return err
	}
	a.Logger = logger.LogDefault
	return nil
}

func (a *App) InitDB() error {

	db, err := mysql.New("default", a.Config.DB.ToMySQL())

	if err != nil {
		return err
	}

	a.DB = db
	return nil
}

func (a *App) InitRedis() error {
	redisClient, err := redis.New(a.Config.Redis.ToRedis(), context.Background())
	if err != nil {
		return err
	}

	a.Redis = redisClient
	return nil
}

func (a *App) InitCache() error {

	redisCache, err := redis.NewCache(a.Config.Redis.ToRedis(), context.Background())
	if err != nil {
		return err
	}
	cacheObj := cache.NewRedisCache(redisCache.Client, nil)

	a.Cache = cacheObj
	return nil
}

func (a *App) InitCaptcha() error {
	captchaObj, err := captcha.NewCaptcha(a.Redis.Client, a.Config.Captcha.ToCaptcha(), nil)
	if err != nil {
		return err
	}
	a.Captcha = captchaObj
	return nil
}

func (a *App) InitStorage() error {
	storageManager, err := storage.New(a.Config.Storage.ToStorage())
	if err != nil {
		return err
	}

	a.Storage = storageManager.Driver()
	return nil
}
