package zrpc

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
)

var once sync.Once
var loger *logrus.Logger
var config *Config

type Config struct {
	Server *Server
	Reg    *Reg
	Db    *gorm.DB
	Log   *logrus.Logger
}

func init()  {
	once.Do(func() {
		Init()
	})
}

func Init() {
	loger = NewLog("log.default_log").loger
	config = InitConfig()
}

func InitConfig() *Config {
	s, err := getServer()
	if err != nil {
		Log().Error(err)
		panic(err)
	}
	db,err := GetDbBySec("mysql.default_mysql")
	if err != nil {
		Log().Error(err)
		panic(err)
	}
	return &Config{
		Server: s,
		Reg:    s.Reg(),
		Db: db,
		Log: Log(),
	}
}

func GetConfig() *Config {
	return config
}

// GetDb 系统默认Db
func GetDb() *gorm.DB {
	return config.Db
}
func Log() *logrus.Logger {
	return loger
}
func GetLog() *logrus.Logger {
	return Log()
}
func GetConf() *Config {
	return config
}
func GetEtcd() *Reg{
	return config.Reg
}