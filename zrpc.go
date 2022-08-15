package zrpc

import (
	"github.com/knadh/koanf"
	"github.com/sirupsen/logrus"
	"github.com/zusux/zrpc/internal"
	"gorm.io/gorm"
	"log"
	"sync"
)

var once sync.Once
var K *koanf.Koanf
var server  *internal.Server


func Init()  {
	once.Do(func() {
		server = InitConfig()
		K = internal.K
	})
}

func InitConfig() *internal.Server {
	s, err := internal.NewServer()
	if err != nil {
		log.Fatalf("[zrpc] init server error: %v",err)
	}
	return s
}



// GetDb 系统默认Db
func GetDb() *gorm.DB {
	db,err:= server.GetDb()
	if err != nil{
		server.Log().Errorf("[zrpc] GetDb error:%s", err)
	}
	return db
}

// NewDb 获取指定Db
func NewDb(section string) (*gorm.DB,error) {
	return server.GetDbBySec(section)
}

// NewLog 获取日志
func NewLog(section string) *logrus.Logger {
	return server.GetLog(section)
}
// Log 获取日志
func Log() *logrus.Logger {
	return server.Log()
}

func GetConf() *internal.Server {
	return server
}

func NewError(code int, message string) *internal.Error {
	return internal.NewError(code,message)
}