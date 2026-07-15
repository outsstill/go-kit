package gokit

import (
	"context"
	"errors"

	_ "time/tzdata"

	"github.com/outsstill/go-kit/cache"
	"github.com/outsstill/go-kit/captcha"
	"github.com/outsstill/go-kit/config"
	"github.com/outsstill/go-kit/database"
	"github.com/outsstill/go-kit/database/mysql"
	"github.com/outsstill/go-kit/jwt"
	"github.com/outsstill/go-kit/limiter"
	"github.com/outsstill/go-kit/logger"
	"github.com/outsstill/go-kit/redis"
	"github.com/outsstill/go-kit/storage"
	"gorm.io/gorm"
)

type GokitApp struct {
	Config  *config.Config
	DB      database.Database
	Redis   *redis.RedisClient
	Cache   cache.Cache
	Captcha *captcha.Captcha
	Storage storage.IStorage
	Logger  *logger.Logger
	JWT     *jwt.JWT
	Limiter *limiter.Limiter
}

var defaultApp *GokitApp

func App() *GokitApp {
	if defaultApp == nil {
		panic("is not initialized")
	}

	return defaultApp
}

type Component int

const (
	Kit_Logger Component = iota
	Kit_DB
	Kit_Redis
	Kit_Cache
	Kit_Captcha
	Kit_Storage
	Kit_JWT
	Kit_Limiter
)

func (a *GokitApp) Init(cs ...Component) error {
	for _, c := range cs {
		var err error

		switch c {
		case Kit_Logger:
			err = a.InitLogger()
		case Kit_DB:
			err = a.InitDB()
		case Kit_Redis:
			err = a.InitRedis()
		case Kit_Cache:
			err = a.InitCache()
		case Kit_Captcha:
			err = a.InitCaptcha()
		case Kit_Storage:
			err = a.InitStorage()
		case Kit_JWT:
			err = a.InitJWT()
		case Kit_Limiter:
			err = a.InitLimiter()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func New(configPath string, cs ...Component) (*GokitApp, error) {
	if defaultApp != nil {
		return defaultApp, nil
	}

	app := &GokitApp{}

	// config
	if err := app.loadConfig(configPath); err != nil {
		return nil, err
	}

	if err := app.Init(cs...); err != nil {
		return nil, err
	}

	defaultApp = app

	return app, nil
}

func (a *GokitApp) loadConfig(configPath string) error {
	configObj, err := config.LoadConfig(configPath)

	if err != nil {
		return errors.New("gokit LoadConfig error: " + err.Error())
	}
	a.Config = configObj
	return nil
}

func (a *GokitApp) InitLogger() error {
	if a.Logger != nil {
		return nil
	}
	err := logger.Init(a.Config.Logger.ToLoggerConfig())
	if err != nil {
		return errors.New("gokit init logger error: " + err.Error())
	}
	a.Logger = logger.LogDefault
	return nil
}

func (a *GokitApp) InitDB() error {

	if a.DB != nil {
		return nil
	}
	db, err := mysql.New("default", a.Config.DB.ToMySQL())

	if err != nil {
		return errors.New("gokit init db error: " + err.Error())
	}

	a.DB = db
	return nil
}

func (a *GokitApp) InitRedis() error {
	if a.Redis != nil {
		return nil
	}
	redisClient, err := redis.New(a.Config.Redis.ToRedis(), context.Background())
	if err != nil {
		return errors.New("gokit init reids error: " + err.Error())
	}

	a.Redis = redisClient
	return nil
}

func (a *GokitApp) InitCache() error {
	if a.Cache != nil {
		return nil
	}
	redisCache, err := redis.NewCache(a.Config.Redis.ToRedis(), context.Background())
	if err != nil {
		return errors.New("gokit init cache error: " + err.Error())
	}
	cacheObj := cache.NewRedisCache(redisCache.Client, nil)

	a.Cache = cacheObj
	return nil
}

func (a *GokitApp) InitCaptcha() error {
	if a.Captcha != nil {
		return nil
	}
	captchaObj, err := captcha.NewCaptcha(a.Redis.Client, a.Config.Captcha.ToCaptcha(), nil)
	if err != nil {
		return errors.New("gokit init captcha error: " + err.Error())
	}
	a.Captcha = captchaObj
	return nil
}

func (a *GokitApp) InitStorage() error {
	if a.Storage != nil {
		return nil
	}
	storageManager, err := storage.New(a.Config.Storage.ToStorage())
	if err != nil {
		return errors.New("gokit init storage error: " + err.Error())
	}

	a.Storage = storageManager.Driver()
	return nil
}

func (a *GokitApp) InitJWT() error {
	if a.JWT != nil {
		return nil
	}
	j, err := jwt.NewJWT(a.Config.JWT.ToJWT())
	if err != nil {
		return errors.New("gokit JWT storage error: " + err.Error())
	}

	a.JWT = j
	return nil
}

func (a *GokitApp) InitLimiter() error {
	if a.Limiter != nil {
		return nil
	}

	if a.Config == nil {
		return errors.New("gokit init limiter config is nil")
	}

	l := limiter.NewLimiter(a.Redis.Client, *a.Config)

	a.Limiter = l
	return nil
}

func Database() database.Database {
	return App().DB
}

func DB() *gorm.DB {
	return App().DB.DB()
}

func Redis() *redis.RedisClient {
	return App().Redis
}

func Config() *config.Config {
	return App().Config
}

func JWT() *jwt.JWT {
	return App().JWT
}

func Cache() cache.Cache {
	return App().Cache
}

func Captcha() *captcha.Captcha {
	return App().Captcha
}

func Log() *logger.Logger {
	return App().Logger
}

func Storage() storage.IStorage {
	return App().Storage
}

func Limiter() *limiter.Limiter {
	return App().Limiter
}
