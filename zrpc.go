package zrpc

import (
	"github.com/knadh/koanf"
	"github.com/sirupsen/logrus"
	"github.com/zusux/zrpc/internal"
	"gorm.io/gorm"
	"sync"
)
var once sync.Once
var config *Config
var NewError  = internal.NewError
var Log func() *logrus.Logger
var K *koanf.Koanf
type Config struct {
	Server *internal.Server
	Reg    *internal.Reg
	Db    *gorm.DB
	Log   *logrus.Logger
}

func Init()  {
	once.Do(func() {
		config = InitConfig()
		K = internal.K
		Log = internal.Log
	})
}

func InitConfig() *Config {
	internal.LoadEnv()
	s, err := internal.GetServer()
	if err != nil {
		internal.Log().Error(err)
		panic(err)
	}
	db,err := internal.GetDbBySec("mysql.default_mysql")
	if err != nil {
		internal.Log().Error(err)
		panic(err)
	}
	return &Config{
		Log: internal.Log(),
		Server: s,
		Reg:    s.Reg(),
		Db: db,
	}
}

func GetConfig() *Config {
	return config
}

// GetDb 系统默认Db
func GetDb() *gorm.DB {
	return config.Db
}

func GetLog() *logrus.Logger {
	return internal.Log()
}
func GetConf() *Config {
	return config
}
func GetEtcd() *internal.Reg{
	return config.Reg
}