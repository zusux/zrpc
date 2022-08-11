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
var Log  func() *logrus.Logger
var GetLog  func() *logrus.Logger
var K *koanf.Koanf
var NewDb func(section string) (*gorm.DB, error)
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
		Log  = internal.Log
		GetLog = internal.Log
		NewDb = internal.GetDbBySec
	})
}

func InitConfig() *Config {
	internal.LoadEnv()
	s, err := internal.GetServer()
	if err != nil {
		internal.Log().Error(err)
	}
	db,err := internal.GetDbBySec("mysql.default_mysql")
	if err != nil {
		GetLog().Error(err)
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

func GetConf() *Config {
	return config
}
func GetEtcd() *internal.Reg{
	return config.Reg
}