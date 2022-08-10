package zrpc

import (
	"github.com/sirupsen/logrus"
	"github.com/zusux/zrpc/internal"
	"gorm.io/gorm"
	"sync"
)
var once sync.Once
var config *Config
var NewError  = internal.NewError
type Config struct {
	Server *internal.Server
	Reg    *internal.Reg
	Db    *gorm.DB
	Log   *logrus.Logger
}

func Init()  {
	once.Do(func() {
		config = InitConfig()
	})
}

func InitConfig() *Config {
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
		Server: s,
		Reg:    s.Reg(),
		Db: db,
		Log: internal.Log(),
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